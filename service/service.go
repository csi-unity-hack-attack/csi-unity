package service

import (
	"context"
	gu "github.com/Murray-LIANG/gounity"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/rexray/gocsi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

const (
	// Name is the name of this CSI SP.
	Name = "csi-unity"

	// VendorVersion is the version of this CSP SP.
	VendorVersion = "0.0.0"
)

// Service is a CSI SP and idempotency.Provider.
type Service interface {
	csi.ControllerServer
	csi.IdentityServer
	csi.NodeServer
	BeforeServe(context.Context, *gocsi.StoragePlugin, net.Listener) error
}

type service struct {
	unityClient *gu.Unity
}

// New returns a new Service.
func New() Service {
	return &service{}
}

func (s *service) BeforeServe(
	ctx context.Context, sp *gocsi.StoragePlugin, lis net.Listener) error {
	if s.unityClient == nil {
		//TODO: Not hard-code
		log.Info("Try to initialize unity client.")
		mgmtIp := "10.141.68.200"
		user := "admin"
		password := "Password123!"
		c, err := gu.NewUnity(mgmtIp, user, password, true)
		if err != nil {
			return status.Errorf(codes.FailedPrecondition,
				"unable to create ScaleIO client: %s", err.Error())
		}
		s.unityClient = c
	}
	return nil
}
