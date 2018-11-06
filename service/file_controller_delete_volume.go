package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	log "github.com/sirupsen/logrus"
)

func FileDeleteVolume(
	s *service,
	ctx context.Context,
	backendId string) error {
	ids := strings.Split(backendId, ":")

	if len(ids) != 2 {
		return errors.New(fmt.Sprintf("Invalid backendId of File volume: %s.", backendId))
	}
	jobErr := deleteVolumeByRest(s.RestEndpoint, ids[0])

	return jobErr
}

func deleteVolumeByRest(rest RestEndpoint, fileSystemId string) error {
	url := fmt.Sprintf("/api/instances/storageResource/%s?compact=true&visibility=Engineering&timeout=0",
		fileSystemId)
	status, resp := rest.delete(url)

	if 404 == status {
		// The delete operation MUST be idempotent as spec
		log.Info("File not existed: ", fileSystemId)
		return nil
	}
	if 202 != status {
		log.Error("Request failed: ", status, resp)
		return errors.New("Delete volume failed.")
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(resp))
	jobId := jsonParsed.Path("id").Data().(string)
	_, _, jobErr := waitForRestJob(rest, jobId)
	if jobErr == nil {
		log.Info("File deleted: ", fileSystemId)
	}
	return jobErr
}
