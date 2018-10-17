package service

import (
	"context"
	gu "github.com/Murray-LIANG/gounity"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
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
