# csi-unity

## Build
* [![Build Status](https://travis-ci.com/jicahoo/csi-unity.svg?branch=master)](https://travis-ci.com/jicahoo/csi-unity)

## Design
* Based on https://github.com/rexray/gocsi and https://github.com/Murray-LIANG/gounity

## How to build&run
* `cd $GOPATH/src/github.com/jicahoo/csi-unity`
* `go install`. This command will generate exe file $GOPATH/bin/csi-unity
* `export CSI_ENDPOINT=csi.sock`
* `$GOPATH/bin/csi-unity`. The command will start the csi-unity server.

## How to test
* Test tool: https://github.com/rexray/gocsi/tree/master/csc
* Install the csc
    * `go get github.com/rexray/gocsi`
    * You will find binary `csc` at $GOPATH/bin
* Prerequisite:
    * `export CSI_ENDPOINT=csi.sock`
* Commands:
    * `./csc controller list-volumes`

## Tools used by this project
* Go package dependency tool: https://github.com/golang/dep


## References about Go
### Go dep
* https://gist.github.com/subfuzion/12342599e26f5094e4e2d08e9d4ad50d
* https://blog.boatswain.io/post/manage-go-dependencies-using-dep/
* https://stackoverflow.com/questions/37237036/how-should-i-use-vendor-in-go-1-6
