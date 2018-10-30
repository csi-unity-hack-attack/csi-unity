package service

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func FileNodeUnpublishVolume(
	s *service,
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {
	return nil, nil
}
