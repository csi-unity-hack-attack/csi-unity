package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Murray-LIANG/gounity"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	hardcodeShare = "fs-ht"
)

type Volume_Type int32

const (
	Volume_Type_Unknow Volume_Type = 0
	Volume_Type_Block  Volume_Type = 1
	Volume_Type_File   Volume_Type = 2
)

var (
	// BlockVolumePrefix is a prefix of Block volume's ID.
	BlockVolumePrefix = fmt.Sprintf("%s-block-", Name)
	// FileVolumePrefix is a prefix of File volume's ID.
	FileVolumePrefix = fmt.Sprintf("%s-file-", Name)
)

// getBackendIdAndTypeByVolumeId is used to get id and volume type by volume id
// if it is block type, returned backend id is block_id
// if it is file type, returned backend id is file_id:nfs_share_id
func getBackendIdAndTypeByVolumeId(volumeId string) (string, Volume_Type) {
	if strings.HasPrefix(volumeId, BlockVolumePrefix) {
		return strings.TrimPrefix(volumeId, BlockVolumePrefix), Volume_Type_Block
	} else if strings.HasPrefix(volumeId, FileVolumePrefix) {
		return strings.TrimPrefix(volumeId, FileVolumePrefix), Volume_Type_File
	} else {
		return volumeId, Volume_Type_Unknow
	}
}

func generateBlockVolumeId(backendId string) string {
	return strings.Join([]string{BlockVolumePrefix, backendId}, "")
}

func generateFileVolumeId(backendId string) string {
	return strings.Join([]string{FileVolumePrefix, backendId}, "")
}

func waitForRestJob(rest RestEndpoint, jobId string) (bool, int, error) {
	completed := false
	state := -1
	var jobErr error = nil
	for i := 0; i < 30; i++ {
		time.Sleep(5 * 1000 * 1000 * 1000)
		completed, state, _ = rest.isJobCompleted(jobId)
		if completed {
			log.Info("Completed. state is ", state)
			if state != 4 {
				jobErr = errors.New(fmt.Sprintf("Job %s failed.", jobId))
			}
			break
		}
	}

	if !completed {
		log.Error("Not completed in time.")
		jobErr = errors.New(fmt.Sprintf("Unity job %s not completed in time", jobId))
	}
	return completed, state, jobErr
}

func (s *service) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	volCaps := req.GetVolumeCapabilities()
	mountCap := volCaps[0].GetMount()

	if mountCap != nil {
		log.Info("Request is to create a 'file' volume.")
		fsType := mountCap.FsType
		log.Info("FsType:", fsType)
		fileResp, fileErr := FileCreateVolume(s, ctx, req)
		if fileErr == nil {
			fileResp.Volume.Id = generateFileVolumeId(fileResp.Volume.Id)
		}
		return fileResp, fileErr
	}

	name := req.GetName()
	log.Info("Try to create volume with name: ", name)

	capRange := req.GetCapacityRange()
	minSize := capRange.RequiredBytes
	maxSize := capRange.LimitBytes
	log.Info("Volume size range (bytes) -- min: ", minSize, " max: ", maxSize)

	//Construct the response
	attrs := make(map[string]string)
	attrs["exportPath"] = "unity_io_ip/my_nfs_share"

	vol := &csi.Volume{
		Id:            hardcodeShare,
		CapacityBytes: int64(1 * gib),
		Attributes:    req.GetParameters(), //Ryan, Why use req.GetParameters?
	}

	resp := &csi.CreateVolumeResponse{
		Volume: vol,
	}

	return resp, nil
}

func (s *service) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {
	log.Info("Call DeleteVolume")

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		log.Error("volume ID missing in request.")
		return nil, status.Error(codes.InvalidArgument, "volume ID missing in request")
	}
	volId := req.GetVolumeId()
	log.Info("Try to delete the volume with id: ", volId)
	if err := s.Driver.ValidateControllerServiceRequest(
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid delete volume req: %v", req)
		return nil, err
	}
	backendId, volType := getBackendIdAndTypeByVolumeId(volId)

	deleteResponse := &csi.DeleteVolumeResponse{}
	var err error = nil
	if volType == Volume_Type_File {
		err = FileDeleteVolume(s, ctx, backendId)
	}

	return deleteResponse, err
}

//TODO: enough for hack attack
func (s *service) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	//Make sure that the volume in request can be accessed by the node in request.
	publishInfo := make(map[string]string)
	publishInfo["export_path"] = "some_value"

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

	resp := &csi.ControllerUnpublishVolumeResponse{}

	return resp, nil
}

//TODO: enough for hack attack
func (s *service) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume ID missing in request")
	}
	if req.GetVolumeCapabilities() == nil {
		return nil, status.Error(
			codes.InvalidArgument, "volume capabilities missing in request")
	}

	//In CSI spec, the method MUST be implemented.
	volId := req.GetVolumeId()
	capabilities := req.GetVolumeCapabilities()
	log.Info("ValidateVolumeCapabilities for vol: ", volId)
	for _, capability := range capabilities {
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
	/*
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
	*/
	//For now, based on Ryan's csi-attack only provide ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
		}}, nil

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

func (s *service) controllerProbe(ctx context.Context) error {

	// Check that we have the details needed to login to the Gateway
	if s.opts.Endpoint == "" {
		return status.Error(codes.FailedPrecondition,
			"missing ScaleIO Gateway endpoint")
	}
	if s.opts.User == "" {
		return status.Error(codes.FailedPrecondition,
			"missing ScaleIO MDM user")
	}
	if s.opts.Password == "" {
		return status.Error(codes.FailedPrecondition,
			"missing ScaleIO MDM password")
	}
	if s.opts.SystemName == "" {
		return status.Error(codes.FailedPrecondition,
			"missing ScaleIO system name")
	}

	// Create our ScaleIO API client, if needed
	if s.unityClient == nil {
		log.Info("Try to initialize unity client. Endpoint:", s.opts.Endpoint, ", user:", s.opts.User)
		mgmtIp := s.opts.Endpoint
		user := s.opts.User
		password := s.opts.Password
		c, err := gounity.NewUnity(mgmtIp, user, password, true)
		if err != nil {
			log.Error("Failed to create Unity client.")
			return status.Errorf(codes.FailedPrecondition,
				"unable to create Unity client: %s", err.Error())
		} else {
			log.Info("Create Unity client successfully.")
		}

		s.SetUnityClient(c)
	}

	// TO DO: Authentication
	return nil
}
