package service

import (
	"bufio"
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"os/exec"
	"strings"

	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
)

const (
	drvCfg = "/opt/emc/unity/sdc/bin/drv_cfg"
)

func (s *service) NodeStageVolume(
	ctx context.Context,
	req *csi.NodeStageVolumeRequest) (
	*csi.NodeStageVolumeResponse, error) {

	return nil, nil
}

func (s *service) NodeUnstageVolume(
	ctx context.Context,
	req *csi.NodeUnstageVolumeRequest) (
	*csi.NodeUnstageVolumeResponse, error) {

	return nil, nil
}

func (s *service) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	return nil, nil
}

func (s *service) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	return nil, nil
}

func (s *service) NodeGetId(
	ctx context.Context,
	req *csi.NodeGetIdRequest) (
	*csi.NodeGetIdResponse, error) {

	return nil, nil
}

func (s *service) NodeGetInfo(
	ctx context.Context,
	req *csi.NodeGetInfoRequest) (
	*csi.NodeGetInfoResponse, error) {

	return nil, nil
}

func (s *service) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	return nil, nil
}

func (s *service) nodeProbe(ctx context.Context) error {

	if s.opts.SdcGUID == "" {
		// try to get GUID using `drv_cfg` binary
		if _, err := os.Stat(drvCfg); os.IsNotExist(err) {
			return status.Error(codes.FailedPrecondition,
				"unable to get SDC GUID via config or drv_cfg binary")
		}

		out, err := exec.Command(drvCfg, "--query_guid").CombinedOutput()
		if err != nil {
			return status.Errorf(codes.FailedPrecondition,
				"error getting SDC GUID: %s", err.Error())
		}

		s.opts.SdcGUID = strings.TrimSpace(string(out))
		//log.WithField("guid", s.opts.SdcGUID).Info("set SDC GUID")
	}

	if !kmodLoaded() {
		return status.Error(codes.FailedPrecondition,
			"scini kernel module not loaded")
	}

	// make sure privDir is pre-created
	if _, err := mkdir(s.privDir); err != nil {
		return status.Errorf(codes.Internal,
			"plugin private dir: %s creation error: %s",
			s.privDir, err.Error())
	}

	return nil
}

func kmodLoaded() bool {
	out, err := exec.Command("lsmod").CombinedOutput()
	if err != nil {
		//log.WithError(err).Error("error from lsmod")
		return false
	}

	r := bytes.NewReader(out)
	s := bufio.NewScanner(r)

	for s.Scan() {
		l := s.Text()
		words := strings.Split(l, " ")
		if words[0] == "scini" {
			return true
		}
	}

	return false
}