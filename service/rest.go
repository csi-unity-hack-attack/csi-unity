package service

import (
	"crypto/tls"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Connection struct {
	ip       string
	username string
	password string
	client   *http.Client
	fields   map[string]string
	useMock  bool
	csrf     string
}

func NewConnection(ip, username, password string) *Connection {
	c := &Connection{ip, username, password, nil, make(map[string]string), false, ""}
	return c.init()
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	// if redirected , use the header of the first request
	//log.WithField("request", req).Debug("request redirected.")
	req.Header = via[0].Header
	return nil
}

func transport() *http.Transport {
	return &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
}

func cookieJar() *cookiejar.Jar {
	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, _ := cookiejar.New(&options)
	//if err != nil {
	//	log.Error(err)
	//}
	return jar
}

func (conn *Connection) init() *Connection {
	conn.client = &http.Client{
		Transport: transport(), Jar: cookieJar(), CheckRedirect: redirectPolicyFunc}
	conn.fields["type"] = ""
	conn.getNewCsrfToken()
	return conn
}

func (conn *Connection) do(req *http.Request) (*http.Response, error) {
	logrus.Info("CSRF: ", conn.csrf)
	req.Header.Set("EMC-CSRF-TOKEN", conn.csrf)
	resp, err := conn.client.Do(req)
	//log.WithField("request", req).Debug("send request.")
	if err != nil {
		//log.WithError(err).Error("http request error.")
		logrus.Error(err)
		return nil, err
	}
	//log.WithField("response", resp).Debug("got response.")
	return resp, err
}

var HEADERS map[string]string = map[string]string{
	"Accept":            "application/json",
	"Content-Type":      "application/json",
	"Accept_Language":   "en_US",
	"X-EMC-REST-CLIENT": "true",
	"User-Agent":        "gounity",
}

func (conn *Connection) newRequest(url, body, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		//log.WithError(err).Error("create request error.")
		return nil, err
	}
	req.SetBasicAuth(conn.username, conn.password)
	for k, v := range HEADERS {
		req.Header.Add(k, v)
	}
	return req, err
}

func (conn *Connection) request(resourcePath, body, method string) (*http.Response, error) {
	unityUrl := fmt.Sprintf("https://%s%s", conn.ip, resourcePath)
	logrus.Info("Request URL: ", unityUrl, " Method: ", method)
	req, err := conn.newRequest(unityUrl, body, method)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = conn.do(req)
	if resp.StatusCode == 401 && method != "GET" {
		req, err := conn.newRequest(unityUrl, body, method)
		if err != nil {
			return nil, err
		}
		resp, err = conn.retryWithCsrfToken(req)
	}
	return resp, err
}

func (conn *Connection) updateCsrf(resp *http.Response) {
	newToken := resp.Header.Get("Emc-Csrf-Token")
	if conn.csrf != newToken {
		conn.csrf = newToken
		//log.WithField("csrf-token", conn.csrf).Info("update csrf token.")
	}
}

func (conn *Connection) getNewCsrfToken() {
	unityUrl := "/api/types/user/instances?fields=name"
	resp, err := conn.request(unityUrl, "", "GET")
	if err != nil {
		logrus.Error("Failed to get new csrf token")
	} else {
		conn.updateCsrf(resp)
	}
}

func (conn *Connection) retryWithCsrfToken(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	//log.Info("token invalid, try to get a new token.")
	conn.getNewCsrfToken()
	if err != nil {
		//log.WithError(err).Error("failed to get csrf-token.")
	} else {
		conn.updateCsrf(resp)
		resp, err = conn.do(req)
	}
	return resp, err
}

func getRespBody(resp *http.Response) string {
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read response body.")
	}
	respBody := string(bytes)
	logrus.WithField("body", respBody).Debug(resp)
	return respBody
}

type RestEndpoint interface {
	get(url string) (int, string)
	post(url string, body string) (int, string)
	delete(url string) (int, string)
	isJobCompleted(jobId string) (bool, int, string)
}

func (conn *Connection) get(resourcePath string) (int, string) {
	re, _ := conn.request(resourcePath, "", "GET")
	respStr := getRespBody(re)
	return re.StatusCode, respStr
}

func (conn *Connection) list(resourceName string) map[string]interface{} {
	return nil
}

func (conn *Connection) post(resourcePath string, body string) (int, string) {
	re, _ := conn.request(resourcePath, body, "POST")
	respStr := getRespBody(re)
	return re.StatusCode, respStr
}

func (conn *Connection) delete(resourcePath string) (int, string) {
	re, _ := conn.request(resourcePath, "", "DELETE")
	respStr := getRespBody(re)
	return re.StatusCode, respStr
}

func (conn *Connection) isJobCompleted(jobId string) (bool, int, string) {
	url := fmt.Sprintf("/api/instances/job/%s/?compact=true&visibility=Engineering&fields=id,description,state,progressPct,endTime,submitTime,tasks,messageOut,parametersOut", jobId)
	re, _ := conn.request(url, "", "GET")
	respStr := getRespBody(re)
	jsonParsed, err := gabs.ParseJSON([]byte(respStr))
	if err != nil {
		logrus.Error("Json parsed error", err)
	}
	state := int(jsonParsed.Path("content.state").Data().(float64))
	logrus.Info("Got status: ", state)
	switch state {
	case 4:
		//Completed successfully
		logrus.Info("Job ", jobId, " completed successfully.")
		return true, 4, ""
	case 5:
		//Failed
		logrus.Error("Job ", jobId, " failed. Job reponse body: \n", respStr)
		return true, 5, ""
	default:
		logrus.Info("Job ", jobId, " is not completed. State is ", state)
		return false, state, ""
	}
}

func EncodeUrl(url string) string {
	encodedQuery := strings.Replace(strings.Replace(url, " ", "%20", -1), "\"", "%22", -1)
	return encodedQuery
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	fmt.Println("Fire!")

	mgmtIp := "10.228.49.124"
	userName := "admin"
	password := "Password123!"
	conn := NewConnection(mgmtIp, userName, password)

	//Query users
	status, respStr := conn.get("/api/types/user/instances?fields=name,role")
	logrus.Info("Got response status: ", status)
	logrus.Debug("Got response body: ", respStr)

	jsonParsed, err := gabs.ParseJSON([]byte(respStr))
	if err != nil {
		logrus.Error("Json parsed error", err)
	}
	entries, _ := jsonParsed.S("entries").Children()
	for _, child := range entries {
		content := child.Path("content").Data().(map[string]interface{})
		logrus.Info(content)
		//logrus.Info(content.Path("name").Data().(string))
		logrus.Info(content["name"])
	}

	//Create user
	create_user_body := `{"name":"ocean", "role":"operator", "password":"Password123!"}`
	status, respStr = conn.post("/api/types/user/instances", create_user_body)
	jsonParsed, err = gabs.ParseJSON([]byte(respStr))
	userId := jsonParsed.Path("content.id").Data().(string)
	logrus.Info("user id: ", userId)

	//Delete user
	unityUrl := fmt.Sprintf("/api/instances/user/%s", userId)
	logrus.Info("Delete user at url: ", unityUrl)
	status, respStr = conn.delete(unityUrl)
	logrus.Info("Status for delete: ", status)

}
