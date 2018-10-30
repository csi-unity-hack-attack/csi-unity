package service

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func FileNodeUnstageVolume(
	s *service,
	ctx context.Context,
	req *csi.NodeUnstageVolumeRequest) (
	*csi.NodeUnstageVolumeResponse, error) {
	return nil, nil
}
