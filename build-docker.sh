#!/bin/bash

go build -o _output/csi-unity .
# Remeber to login
docker build -t ciqihuo/csi-unity:0.1 .

