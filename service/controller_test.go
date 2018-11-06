package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	gu "github.com/Murray-LIANG/gounity"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/stretchr/testify/assert"
)

type mockUnity struct {
}

func (m *mockUnity) GetPools() ([]*gu.Pool, error) {
	return nil, nil
}

func (m *mockUnity) GetPoolById(id string) (*gu.Pool, error) {
	return nil, nil
}

func (m *mockUnity) GetPoolByName(name string) (*gu.Pool, error) {
	return nil, nil
}

func (m *mockUnity) GetLuns() ([]*gu.Lun, error) {
	size := 10
	luns := make([]*gu.Lun, size)
	for i := 0; i < size; i++ {
		lun := &gu.Lun{
			Id:            "sv_" + strconv.Itoa(i),
			SizeAllocated: 100,
		}
		luns[i] = lun /* assign the address of integer. */
	}

	return luns, nil
}
func (m *mockUnity) GetLunById(id string) (*gu.Lun, error) {
	return nil, nil
}

func (m *mockUnity) GetLunByName(name string) (*gu.Lun, error) {
	return nil, nil
}

func TestListVolumes(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	m := &mockUnity{}
	ms := &service{}
	ms.SetUnityClient(m)

	ctx := context.Context(context.Background())

	//Case: big max entries
	req := &csi.ListVolumesRequest{MaxEntries: 100, StartingToken: ""}
	resp, _ := ms.ListVolumes(ctx, req)

	assert.Equal(t, 10, len(resp.Entries))
	assert.Equal(t, "", resp.NextToken)

	//Case: little max entries
	req = &csi.ListVolumesRequest{MaxEntries: 2, StartingToken: ""}
	resp, _ = ms.ListVolumes(ctx, req)
	assert.Equal(t, 2, len(resp.Entries))
}

func TestService_GetCapacity(t *testing.T) {
	m := &mockUnity{}
	ms := &service{}
	ms.SetUnityClient(m)

	ctx := context.Context(context.Background())
	req := &csi.GetCapacityRequest{}
	resp, _ := ms.GetCapacity(ctx, req)
	assert.Equal(t, int64(10*1024*1024*1024*1024), resp.AvailableCapacity)
}

func TestService_ControllerGetCapabilities(t *testing.T) {
	ms := &service{}

	ctx := context.Context(context.Background())
	req := &csi.ControllerGetCapabilitiesRequest{}
	resp, _ := ms.ControllerGetCapabilities(ctx, req)
	for _, capability := range resp.Capabilities {
		fmt.Println(capability.String())
		found := false
		capStr := capability.String()
		if "rpc:<type:CREATE_DELETE_VOLUME > " == capStr {
			found = true
		}
		if "rpc:<type:PUBLISH_UNPUBLISH_VOLUME > " == capStr {
			found = true
		}
		if "rpc:<type:LIST_VOLUMES > " == capStr {
			found = true
		}

		if "rpc:<type:GET_CAPACITY > " == capStr {
			found = true
		}
		assert.Equal(t, true, found, "Unexpect capability: ", capStr)
	}

}

func Test_GetBackendIdAndTypeByVolumeId_UnknowType(t *testing.T) {
	exceptedId := "test_volume_id"
	volId := exceptedId
	exceptedVolumeType := Volume_Type_Unknow
	id, volType := getBackendIdAndTypeByVolumeId(volId)
	assert.Equal(t, exceptedId, id, "UnExcepted id: ", id)
	assert.Equal(t, exceptedVolumeType, volType, "UnExcepted volume type: ", volType)
}

func Test_GetBackendIdAndTypeByVolumeId_BlockType(t *testing.T) {
	exceptedId := "test_volume_id"
	volId := BlockVolumePrefix + exceptedId
	exceptedVolumeType := Volume_Type_Block
	id, volType := getBackendIdAndTypeByVolumeId(volId)
	assert.Equal(t, exceptedId, id, "UnExcepted id: ", id)
	assert.Equal(t, exceptedVolumeType, volType, "UnExcepted volume type: ", volType)
}

func Test_GetBackendIdAndTypeByVolumeId_FileType(t *testing.T) {
	exceptedId := "test_volume_id:share_id"
	volId := FileVolumePrefix + exceptedId
	exceptedVolumeType := Volume_Type_File
	id, volType := getBackendIdAndTypeByVolumeId(volId)
	assert.Equal(t, exceptedId, id, "UnExcepted id: ", id)
	assert.Equal(t, exceptedVolumeType, volType, "UnExcepted volume type: ", volType)
}

func Test_GenerateBlockVolumeId(t *testing.T) {
	id := "test_id"
	exceptedVolId := BlockVolumePrefix + id
	volId := generateBlockVolumeId(id)
	assert.Equal(t, exceptedVolId, volId, "UnExcepted volume id: ", volId)
}

func Test_GenerateFileVolumeId(t *testing.T) {
	id := "test_id"
	exceptedVolId := FileVolumePrefix + id
	volId := generateFileVolumeId(id)
	assert.Equal(t, exceptedVolId, volId, "UnExcepted volume id: ", volId)
}
