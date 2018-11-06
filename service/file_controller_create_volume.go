package service

import (
	"bytes"
	"container/list"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/Jeffail/gabs"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/sirupsen/logrus"
)

type fileInterfaceInput struct {
	nasServer string
	ipPort string
	ipAddress string
	netmask string
	gateway string
}

func FileCreateVolume(
	s *service,
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {
	return nil, nil

	name := req.GetName()
	capRange := req.GetCapacityRange()
	size := capRange.GetRequiredBytes()

    //create io interface for volume first
    ipMeta :=fileInterfaceInput{"nas_1", "spa_iom_0_eth0", "10.103.76.148","255.255.248.0", "10.103.72.1"}

	_, err := createFileInterface(s.RestEndpoint, ipMeta)
	if err != nil {
		logrus.Error("Error: Failed to create file interface.")
		//ignore failure in file interface creation
		//return nil, err
	}

	volAttrMap, _ := createVolumeByRest(s.RestEndpoint, uint64(size), name)
	vol := &csi.Volume{
		Id:            volAttrMap["storageResource"] + ":" + volAttrMap["id"],
		CapacityBytes: size,
		Attributes:    volAttrMap,
	}

	resp := &csi.CreateVolumeResponse{
		Volume: vol,
	}
	return resp, nil
}

func createVolumeByRest(rest RestEndpoint, size uint64, name string) (map[string]string, error) {
	paras := map[string]string{
		"FsName":    name,
		"FsSize":    strconv.Itoa(int(size)),
		"ShareName": name,
		"PoolId":    "pool_1",
		"NasId":     "nas_1",
	}

	requestTemplate := `{
    "description": "Creating File System",
    "tasks": [
        {
            "name": "CreateNewFilesystem",
            "object": "storageResource",
            "action": "createFilesystem",
            "description": "Create File System",
            "parametersIn": {
                "name": "{{.FsName}}",
                "description": "",
                "fsParameters": {
                    "pool": {
                        "id": "{{.PoolId}}"
                    },
                    "nasServer": {
                        "id": "{{.NasId}}"
                    },
                    "isThinEnabled": true,
                    "supportedProtocols": 0,
                    "format": 2,
                    "size": {{.FsSize}},
                    "isDataReductionEnabled": false,
                    "isAdvancedDedupEnabled": false,
                    "fastVPParameters": {
                        "tieringPolicy": 0
                    },
                    "flrVersion": 0
                },
                "cifsFsParameters": {
                    "isCIFSOpLocksEnabled": false
                },
                "nfsShareCreate": [
                    {
                        "name": "{{.ShareName}}",
                        "path": "/",
                        "nfsShareParameters": {
                            "description": "",
                            "defaultAccess": 0,
                            "noAccessHosts": [],
                            "readOnlyHosts": [],
                            "readWriteHosts": [],
                            "rootAccessHosts": [],
                            "readOnlyRootAccessHosts": [],
                            "anonymousUID": 4294967294,
                            "anonymousGID": 4294967294
                        }
                    }
                ]
            }
        }
    ],
    "majorStepIndex": 1}`
	result := populateParas(requestTemplate, paras)
	fmt.Println(result)

	url := "/api/types/job/instances?compact=true&visibility=Engineering&timeout=0"
	status, resp := rest.post(url, result)
	if 202 != status {
		logrus.Error("Request failed: ", status, resp)
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(resp))
	jobId := jsonParsed.Path("id").Data().(string)

	completed, _, jobErr := waitForRestJob(rest, jobId)

	if completed {
		return queryShareData(rest, name)
	}

	return nil, jobErr
}

func populateParas(reqTemplate string, paras map[string]string) string {

	tmpl, err := template.New("request").Parse(reqTemplate)
	if err != nil {
		logrus.Error("Request template render error: ", err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, paras)
	if err != nil {
		logrus.Error("Request template render error: ", err)
	}
	result := tpl.String()
	hasNoValue := strings.Contains(result, "<no value>")
	if hasNoValue {
		logrus.Error("Found <no value> in rendered string: ", result)
	}
	return result
}

func queryShareData(conn RestEndpoint, fsName string) (map[string]string, error) {
	unityUrl := fmt.Sprintf(`/api/types/nfsShare/instances?with_entrycount=true&compact=true&visibility=Engineering&fields=name,type,nasServerName::filesystem.nasServer.name,filesystemName::filesystem.name,localPath::((snap ne null) ? @concat("/", snap.name, path) : @concat("/", filesystem.name, path)),id,description,exportPath::@concatList(filesystem.nasServer.fileInterface.ipAddress),snap.name,role,isReadOnly,hostAccessesCount::@count(hostAccesses),creationTime,modificationTime,defaultAccess,storageResource::filesystem.storageResource.id,interfacesCount::@count(filesystem.nasServer.fileInterface),filesystem.nasServer.isReplicationDestination&page=1&per_page=100&orderby=name ASC&filter=(filesystem.type eq 1) and (filesystem.name eq "%s")`, fsName)
	encodedQuery := EncodeUrl(unityUrl)
	logrus.Debug("Encoded: ", encodedQuery)
	status, resp := conn.get(encodedQuery)
	logrus.Debug("Status: ", status)
	logrus.Debug("Resp: ", resp)
	jsonParsed, _ := gabs.ParseJSON([]byte(resp))
	entryCount := jsonParsed.Path("entryCount").Data().(float64)
	nfsShareMetaData := make(map[string]string)

	var err error = nil
	if int(entryCount) != 1 {
		logrus.Error("Error, the number of nfs share entries is not 1 for fs name: ", fsName)
		err = errors.New(fmt.Sprintf("Error, the number of nfs share entries is not 1 for fs name: %s", fsName))
	} else {
		children, _ := jsonParsed.S("entries").Children()
		for _, child := range children {
			content := child.Path("content")
			nfsShareMetaData["exportPath"] = content.Path("exportPath").Data().(string)
			nfsShareMetaData["localPath"] = content.Path("localPath").Data().(string)
			nfsShareMetaData["nasServerName"] = content.Path("nasServerName").Data().(string)
			nfsShareMetaData["name"] = content.Path("name").Data().(string)
			nfsShareMetaData["id"] = content.Path("id").Data().(string)
			nfsShareMetaData["storageResource"] = content.Path("storageResource").Data().(string)
			nfsShareMetaData["filesystemName"] = content.Path("filesystemName").Data().(string)
		}
		logrus.Info(nfsShareMetaData)
	}
	return nfsShareMetaData, err
}

func queryFileInterface(conn RestEndpoint, nasServer string) (*list.List, error){
	var getFileInterfaceUrl string = "/api/types/fileInterface/instances?fields=id,nasServer,ipAddress"
    encodedQuery := EncodeUrl(getFileInterfaceUrl)
    logrus.Debug("Encoded: ", encodedQuery)
    status, resp := conn.get(encodedQuery)
    logrus.Debug("Status: ", status)
    logrus.Debug("Query file interface response: ", resp)
    jsonParsed, _ := gabs.ParseJSON([]byte(resp))
    entryCount := jsonParsed.Path("entryCount").Data().(float64)
    fileInterfaces := list.New();


    var err error = nil
    if int(entryCount) == 0 {
    	logrus.Error("Error: the num of file interface entries is 0")
		err = errors.New(fmt.Sprintf("Error, tError: the num of file interface entries is 0 for nas server: %s", nasServer))
	}else{
		children, _ := jsonParsed.S("entries").Children()
		for _, child := range children{
			if child.Path("content").Path("nasServer").Data().(string) == nasServer {
				fileInterfaceData := make(map[string]string)
				fileInterfaceData["id"] = child.Path("content").Path("id").Data().(string)
				fileInterfaceData["ipAddr"] = child.Path("content").Path("ipAddress").Data().(string)
				fileInterfaces.PushBack(fileInterfaceData);
				logrus.Info(fileInterfaceData)
			}
		}
		logrus.Info(fileInterfaces)
	}
    return fileInterfaces, err
}

func createFileInterface(conn RestEndpoint, fileIPInput fileInterfaceInput)(*list.List, error) {

	fileInterfaces, err := queryFileInterface(conn, fileIPInput.nasServer)
	if fileInterfaces.Len() == 0 {
		logrus.Info("Info: file interface is not created for nas server %s, start to create now.", fileIPInput.nasServer)

		// start to create interface
		var createFileInterfaceUrl string = "/api/types/fileInterface/instances?compact=true&visibility=Engineering&timeout=0"
		// Body: {"nasServer": {"id":"nas_1"}, "ipPort": {"id":"spa_iom_0_eth0"}, "ipAddress": "10.103.76.143","netmask":"255.255.248.0", "gateway":"10.103.72.1"}
		var createFileInterfaceBody string = fmt.Sprintf(`{"nasServer": {"id":"%s"}, "ipPort": {"id":"%s"}, "ipAddress":"%s", "netmask":"%s", "gateway":"%s"}`, fileIPInput.nasServer, fileIPInput.ipPort, fileIPInput.ipAddress, fileIPInput.netmask, fileIPInput.gateway)
		status, resp := conn.post(createFileInterfaceUrl, createFileInterfaceBody)

		logrus.Debug("Status: ", status)
		logrus.Debug("Create file interface response: ", resp)

		if 202 != status {
			logrus.Error("Request failed: ", status, resp)
		}

		jsonParsed, _ := gabs.ParseJSON([]byte(resp))
		jobId := jsonParsed.Path("id").Data().(string)

		completed, _, jobErr := waitForRestJob(conn, jobId)

		if completed {
			return queryFileInterface(conn, fileIPInput.nasServer)
		}
		return nil, jobErr

	} else {
		logrus.Info("File interface(s) are already created for nas server %s", fileIPInput.nasServer)
		return fileInterfaces,err
	}
}
