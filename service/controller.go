package service

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"golang.org/x/net/context"
)

func (s *service) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	return nil, nil
}

func (s *service) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	return nil, nil
}

func (s *service) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	return nil, nil
}

func (s *service) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	return nil, nil
}

func (s *service) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	return nil, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (s *service) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	//Extract info from request
	if v := req.StartingToken; v != "" {
		//TODO: support
	}
	maxEntries := req.MaxEntries

	luns, _ := s.unityClient.GetLuns()
	lenOfLuns := len(luns)
	respSize := min(int(maxEntries), lenOfLuns)

	entries := make(
		[]*csi.ListVolumesResponse_Entry,
		respSize)
	nextToken := ""
	for i := 0; i < int(respSize); i++ {
		vi := &csi.Volume{
			Id:            luns[i].Id,
			CapacityBytes: (int64)(luns[i].SizeTotal),
		}
		//TODO: support: nextToken = luns[i].Id
		entries[i] = &csi.ListVolumesResponse_Entry{
			Volume: vi,
		}
	}

	return &csi.ListVolumesResponse{
		Entries:   entries,
		NextToken: nextToken,
	}, nil
}

func (s *service) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {
	tenTb := 10 * 1024 * 1024 * 1024 * 1024
	return &csi.GetCapacityResponse{AvailableCapacity: int64(tenTb)}, nil
}

func (s *service) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	return nil, nil
}

func (s *service) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest) (
	*csi.CreateSnapshotResponse, error) {

	return nil, nil
}

func (s *service) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest) (
	*csi.DeleteSnapshotResponse, error) {

	return nil, nil
}

func (s *service) ListSnapshots(
	ctx context.Context,
	req *csi.ListSnapshotsRequest) (
	*csi.ListSnapshotsResponse, error) {

	return nil, nil
}
