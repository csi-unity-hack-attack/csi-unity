package service

import (
	"context"
	"strings"
	"testing"

	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRest struct {
	mock.Mock
}

func (m *mockRest) delete(resourcePath string) (int, string) {
	if strings.Contains(resourcePath, "/not_existed_id") {
		return 404, ""
	} else if strings.Contains(resourcePath, "/error_id") {
		return 400, ""
	} else {
		return 202, "{\"id\": \"1\"}"
	}
}

func (m *mockRest) isJobCompleted(jobId string) (bool, int, string) {
	switch jobId {
	case "1":
		return true, 4, ""
	default:
		return true, 4, ""
	}
}

func (m *mockRest) get(url string) (int, string) {
	return 200, ""
}
func (m *mockRest) post(url string, body string) (int, string) {
	return 200, ""
}

var restEndpoint = &mockRest{}
var driver = csicommon.NewCSIDriver("test", "0.0.1", "ubuntu")
var svr = &service{
	RestEndpoint: restEndpoint,
	Driver:       driver,
}

func configDriver() {
	driver.AddControllerServiceCapabilities(
		[]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		})
}

func TestDeleteVolumeByFile(t *testing.T) {
	configDriver()
	req := &csi.DeleteVolumeRequest{
		VolumeId: "csi-unity-file-test_file_id:test_share_id",
	}
	_, err := svr.DeleteVolume(context.Background(), req)
	assert.Empty(t, err, "Excepted no error, but got error: ", err)
}

func TestDeleteVolumeNotExistedByFile(t *testing.T) {
	configDriver()
	req := &csi.DeleteVolumeRequest{
		VolumeId: "csi-unity-file-not_existed_id:test_share_id",
	}
	_, err := svr.DeleteVolume(context.Background(), req)
	assert.Empty(t, err, "Excepted no error, but got error: ", err)
}

func TestDeleteVolumeErrorByFile(t *testing.T) {
	configDriver()
	req := &csi.DeleteVolumeRequest{
		VolumeId: "csi-unity-file-error_id:test_share_id",
	}
	_, err := svr.DeleteVolume(context.Background(), req)
	assert.NotEmpty(t, err, "Excepted error, but no error")
}

func TestFileDeleteVolume(t *testing.T) {
	err := FileDeleteVolume(svr, context.Background(), "test_file_id:test_share_id")
	assert.Empty(t, err, "Excepted no error, but got error: ", err)
}

func TestFileDeleteVolumeNotExisted(t *testing.T) {
	err := FileDeleteVolume(svr, context.Background(), "not_existed_id:test_share_id")
	assert.Empty(t, err, "Excepted no error, but got error: ", err)
}

func TestFileDeleteVolumeError(t *testing.T) {
	err := FileDeleteVolume(svr, context.Background(), "error_id:test_share_id")
	assert.NotEmpty(t, err, "Excepted error, but np error")
}
