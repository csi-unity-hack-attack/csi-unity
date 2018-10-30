package service

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func FileNodeStageVolume(
	s *service,
	ctx context.Context,
	req *csi.NodeStageVolumeRequest) (
	*csi.NodeStageVolumeResponse, error) {
	return nil, nil
}
