package service

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func FileControllerPublishVolume(
	s *service,
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {
	return nil, nil
}
