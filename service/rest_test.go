package service

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestRestCrud(t *testing.T) {
	unityIp, hasEnv := os.LookupEnv(UtUnityIp)

	if hasEnv {
		//unityIp := "10.228.49.124"
		userName := "admin"
		password := "Password123!"
		conn := NewConnection(unityIp, userName, password)

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

}
