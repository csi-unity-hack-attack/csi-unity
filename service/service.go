package service

import (
	"context"
	gu "github.com/Murray-LIANG/gounity"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/rexray/gocsi"
	csictx "github.com/rexray/gocsi/context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

const (
	// Name is the name of this CSI SP.
	Name = "csi-unity"

	// VendorVersion is the version of this CSP SP.
	VendorVersion  = "0.0.0"
	defaultPrivDir = "/dev/disk/csi-unity"
)

// Manifest is the SP's manifest.
var Manifest = map[string]string{
	"url":    "https://github.com/thecodeteam/csi-scaleio",
	"semver": "1.0.0",
}

// Service is a CSI SP and idempotency.Provider.
type Service interface {
	csi.ControllerServer
	csi.IdentityServer
	csi.NodeServer
	BeforeServe(context.Context, *gocsi.StoragePlugin, net.Listener) error
}

type service struct {
	Driver       *csicommon.CSIDriver
	unityClient  gu.Storage
	mode         string
	opts         Opts
	privDir      string
	NodeId       string
	RestEndpoint RestEndpoint
}

// New returns a new Service.
func New() Service {
	return &service{}
}

func (s *service) SetUnityClient(storage gu.Storage) {
	s.unityClient = storage
}

type Opts struct {
	Endpoint   string
	User       string
	Password   string
	SystemName string
	SdcGUID    string
	Insecure   bool
	Thick      bool
	AutoProbe  bool
	NodeId     string
}

func (s *service) BeforeServe(
	ctx context.Context, sp *gocsi.StoragePlugin, lis net.Listener) error {

	s.mode = csictx.Getenv(ctx, gocsi.EnvVarMode)
	opts := Opts{}
	s.opts = opts

	if ep, ok := csictx.LookupEnv(ctx, Endpoint); ok {
		opts.Endpoint = ep
	}
	if user, ok := csictx.LookupEnv(ctx, User); ok {
		opts.User = user
	}
	if opts.User == "" {
		opts.User = "admin"
	}
	if pw, ok := csictx.LookupEnv(ctx, Password); ok {
		opts.Password = pw
	}

	if pd, ok := csictx.LookupEnv(ctx, "X_CSI_PRIVATE_MOUNT_DIR"); ok {
		s.privDir = pd
	}

	if "" == s.privDir {
		s.privDir = defaultPrivDir
	}

	if nodeId, ok := csictx.LookupEnv(ctx, NodeId); ok {
		opts.NodeId = nodeId
		s.NodeId = nodeId
		log.Info("Get node id: ", nodeId)
	}

	//TODO: refine -- Prepare driver.
	s.Driver = csicommon.NewCSIDriver(Name, VendorVersion, s.NodeId)
	if s.Driver == nil {
		glog.Fatalln("failed to initialize csi driver.")
	} else {
		log.Info("Got csicommon CSIDriver.")
	}
	s.Driver.AddControllerServiceCapabilities(
		[]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			// TODO (ryan) add snapshots related capabilities
		})
	s.Driver.AddVolumeCapabilityAccessModes(
		[]csi.VolumeCapability_AccessMode_Mode{
			csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		})

	if s.unityClient == nil {
		log.Info("Try to initialize unity client. Endpoint:", opts.Endpoint, ", user:", opts.User)
		mgmtIp := opts.Endpoint
		user := opts.User
		password := opts.Password
		c, err := gu.NewUnity(mgmtIp, user, password, true)
		if err != nil {
			log.Error("Failed to create Unity client.")
			return status.Errorf(codes.FailedPrecondition,
				"unable to create Unity client: %s", err.Error())
		} else {
			log.Info("Create Unity client successfully.")
		}

		s.SetUnityClient(c)
	}

	if s.RestEndpoint == nil {
		log.Info("Try to initialize the REST endpoint.")
		conn := NewConnection(opts.Endpoint, opts.User, opts.Password)
		s.RestEndpoint = conn
	}
	return nil
}
