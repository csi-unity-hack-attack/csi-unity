package service

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func FileNodePublishVolume(
	s *service,
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {
	return nil, nil
}
