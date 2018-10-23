package service

import (
	"fmt"
	"os"
)

// Device is a struct for holding details about a block device
type Device struct {
	FullPath string
	Name     string
	RealDev  string
}



// mkdir creates the directory specified by path if needed.
// return pair is a bool flag of whether dir was created, and an error
func mkdir(path string) (bool, error) {
	st, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			//log.WithField("dir", path).WithError(
			//	err).Error("Unable to create dir")
			return false, err
		}
		//log.WithField("path", path).Debug("created directory")
		return true, nil
	}
	if !st.IsDir() {
		return false, fmt.Errorf("existing path is not a directory")
	}
	return false, nil
}