package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileCreateVolume(t *testing.T) {
	//TODO
	unityIp, hasEnv := os.LookupEnv(UtUnityIp)
	if hasEnv {
		fmt.Println("csi-unity-002")
		//unityIp := "10.228.49.124"
		userName := "admin"
		password := "Password123!"
		conn := NewConnection(unityIp, userName, password)
		_, jobErr := createVolumeByRest(conn, uint64(10*gib), "csi-unity-002")
		assert.True(t, jobErr == nil, jobErr.Error())
	}
}
