package service

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//TODO: not enough for hack attack
func (s *service) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {



	name := req.GetName()
	log.Info("Try to create volume with name: ", name)

	capRange := req.GetCapacityRange()
	minSize := capRange.RequiredBytes
	maxSize := capRange.LimitBytes
	log.Info("Volume size range (bytes) -- min: ", minSize, " max: ", maxSize )

	//TODO: Call to Unity

	//Construct the response
	attrs := make(map[string]string)
	attrs["exportPath"] = "unity_io_ip/my_nfs_share"

	vol := &csi.Volume{
		Id:     "nfs_1"       ,
		CapacityBytes: 0, //0 for nfs
		Attributes: attrs,
	}

	resp := &csi.CreateVolumeResponse{
		Volume: vol,
	}

	return resp, nil
}

//TODO: not enough for hack attack
func (s *service) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	volId := req.GetVolumeId()
	log.Info("Try to delete the volume with id: ", volId)
	//TODO: determine if it is a NFS or LUN. Then send request to Unity

	return &csi.DeleteVolumeResponse{}, nil
}

//TODO: enough for hack attack
func (s *service) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	//Make sure that the volume in request can be accessed by the node in request.
	publishInfo := make(map[string]string)
	publishInfo["some_key"] = "some_value"

	resp := &csi.ControllerPublishVolumeResponse{
		PublishInfo: publishInfo,
	}
	return resp, nil
}

//TODO: enough for hack attack
func (s *service) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	//Unpublish the volume form all nodes or node specified in request.

	//TODO: implement the function

	resp := &csi.ControllerUnpublishVolumeResponse{

	}

	return resp, nil
}

//TODO: enough for hack attack
func (s *service) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	//In CSI spec, the method MUST be implemented.
	volId := req.GetVolumeId()
	capabilities := req.GetVolumeCapabilities()
	log.Info("ValidateVolumeCapabilities for vol: ", volId)
	for _,capability := range capabilities {
		log.Info("Capability access mode is: ", capability.GetAccessMode().GetMode().String())
		log.Info("Capability access type is: ", capability.GetAccessType())
	}

	//TODO: For now, return true for all case
	resp := &csi.ValidateVolumeCapabilitiesResponse{
		Supported: true,
	}

	return resp, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

//TODO: Not enough for hack attack
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

//TODO: enough for hack attack
func (s *service) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {
	//TODO: I believe such implementation is enough for Demo.
	tenTb := 10 * 1024 * 1024 * 1024 * 1024
	return &csi.GetCapacityResponse{AvailableCapacity: int64(tenTb)}, nil
}

//TODO: enough for hack attack
func (s *service) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_GET_CAPACITY,
					},
				},
			},
		},
	}, nil
}

//TODO: enough for hack attack
func (s *service) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest) (
	*csi.CreateSnapshotResponse, error) {

	//TODO: Create snapshot from volume, and return the created snap info.
	return nil, status.Error(codes.Unimplemented, "")

}

//TODO: enough for hack attack
func (s *service) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest) (
	*csi.DeleteSnapshotResponse, error) {

	return nil, status.Error(codes.Unimplemented, "")
}


//TODO: enough for hack attack
func (s *service) ListSnapshots(
	ctx context.Context,
	req *csi.ListSnapshotsRequest) (
	*csi.ListSnapshotsResponse, error) {
	//TODO: return snapshots according to volume id or snapshot id.
	return nil, status.Error(codes.Unimplemented, "")
}
