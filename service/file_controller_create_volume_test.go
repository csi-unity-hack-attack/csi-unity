package service

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileCreateVolume(t *testing.T) {
	//TODO
	unityIp, hasEnv := os.LookupEnv(UtUnityIp)
	if hasEnv {
		//unityIp := "10.228.49.124"
		userName := "admin"
		password := "Password123!"
		conn := NewConnection(unityIp, userName, password)
		volumeName := "csi-unity-003"
		nfsShareData, jobErr := createVolumeByRest(conn, uint64(10*gib), volumeName)
		assert.True(t, jobErr == nil, "Job err is not nil")
		logrus.Info(nfsShareData)
	}
}
