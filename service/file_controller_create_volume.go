package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func FileCreateVolume(
	s *service,
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {
	return nil, nil

	name := req.GetName()
	capRange := req.GetCapacityRange()
	size := capRange.GetRequiredBytes()
	volAttrMap, _ := createVolumeByRest(s.RestEndpoint, uint64(size), name)
	vol := &csi.Volume{
		Id:            volAttrMap["id"],
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

	completed := false
	state := -1
	var jobErr error = nil
	for i := 0; i < 30; i++ {
		time.Sleep(5 * 1000 * 1000 * 1000)
		completed, state, _ = rest.isJobCompleted(jobId)
		if completed {
			logrus.Info("Completed. state is ", state)
			if state != 4 {
				jobErr = errors.New(fmt.Sprintf("Job %s failed.", jobId))
			}
			break
		}

	}

	if !completed {
		logrus.Error("Not completed in time.")
	} else {
		//TODO: get nfs share info direclty by name: https://10.228.49.124/api/types/nfsShare/instances?fields=name&filter=name%20eq%20%22csi-unity-002%22
		url = `/api/types/nfsShare/instances?fields=name&filter=name eq "csi-unity-002"`
		//GUI URL: https://10.228.49.124/api/instances/nfsShare/NFSShare_4/?with_entrycount=true&compact=true&visibility=Engineering&fields=name%2Ctype%2CnasServerName%3A%3Afilesystem.nasServer.name%2CfilesystemName%3A%3Afilesystem.name%2ClocalPath%3A%3A((snap%20ne%20null)%20%3F%20%40concat(%22%2F%22%2C%20snap.name%2C%20path)%20%3A%20%40concat(%22%2F%22%2C%20filesystem.name%2C%20path))%2Cid%2Cdescription%2CexportPath%3A%3A%40concatList(filesystem.nasServer.fileInterface.ipAddress)%2Csnap.name%2Crole%2CisReadOnly%2ChostAccessesCount%3A%3A%40count(hostAccesses)%2CcreationTime%2CmodificationTime%2CdefaultAccess%2CstorageResource%3A%3Afilesystem.storageResource.id%2CinterfacesCount%3A%3A%40count(filesystem.nasServer.fileInterface)%2Cfilesystem.nasSer
		// ver.isReplicationDestination&page=1&per_page=100&orderby=name%20ASC&filter=(filesystem.type%20eq%201)
		//https://10.228.49.124/api/instances/nfsShare/NFSShare_4/?with_entrycount=true&compact=true&visibility=Engineering&fields=name,type,nasServerName::filesystem.nasServer.name,filesystemName::filesystem.name,localPath::((snap ne null) ? @concat("/", snap.name, path) : @concat("/", filesystem.name, path)),id,description,exportPath::@concatList(filesystem.nasServer.fileInterface.ipAddress),snap.name,role,isReadOnly,hostAccessesCount::@count(hostAccesses),creationTime,modificationTime,defaultAccess,storageResource::filesystem.storageResource.id,interfacesCount::@count(filesystem.nasServer.fileInterface),filesystem.nasServer.isReplicationDestination&page=1&per_page=100&orderby=name ASC&filter=(filesystem.type eq 1)
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
