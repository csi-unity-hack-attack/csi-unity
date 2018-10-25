#!/bin/bash
export CSI_ENDPOINT=tcp://127.0.0.1:34555
#mgmt IP
export X_CSI_UNITY_ENDPOINT=10.228.49.124
export X_CSI_UNITY_USER=admin
export X_CSI_UNITY_PASSWORD=Password123!
export KUBE_NODE_NAME=node-1
go run main.go
