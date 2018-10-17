#!/bin/bash
export CSI_ENDPOINT=tcp://127.0.0.1:34555
#mgmt IP
export CSI_X_UNITY_ENDPOINT=10.141.68.200
export CSI_X_UNITY_USER=admin
export CSI_X_UNITY_PASSWORD=<YourPassword>
go run main.go
