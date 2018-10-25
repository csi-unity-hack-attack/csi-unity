package service

import (
	"golang.org/x/net/context"
	"strings"

	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
)

func (s *service) Probe(
	ctx context.Context,
	req *csi.ProbeRequest) (
	*csi.ProbeResponse, error) {
	if strings.EqualFold(s.mode, "both") {
		return &csi.ProbeResponse{}, nil
	} else {
		if strings.EqualFold(s.mode, "controller") {
			if err := s.controllerProbe(ctx); err != nil {
				return nil, err
			}
		}
		if strings.EqualFold(s.mode, "node") {
			if err := s.nodeProbe(ctx); err != nil {
				return nil, err
			}
		}
	}
	return &csi.ProbeResponse{}, nil

}

func (s *service) GetPluginInfo(
	ctx context.Context,
	req *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {

	return &csi.GetPluginInfoResponse{
		Name:          Name,
		VendorVersion: VendorVersion,
		Manifest:      Manifest,
	}, nil
}

func (s *service) GetPluginCapabilities(
	ctx context.Context,
	req *csi.GetPluginCapabilitiesRequest) (
	*csi.GetPluginCapabilitiesResponse, error) {

	var rep csi.GetPluginCapabilitiesResponse
	if !strings.EqualFold(s.mode, "node") {
		rep.Capabilities = []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		}
	}
	return &rep, nil
}
