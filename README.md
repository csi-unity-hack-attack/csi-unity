# csi-unity
## Project Status
* In developing!

## Build
* [![Build Status](https://travis-ci.com/jicahoo/csi-unity.svg?branch=master)](https://travis-ci.com/jicahoo/csi-unity)

## Design
* Based on https://github.com/rexray/gocsi and https://github.com/Murray-LIANG/gounity

## Contribute Flow
* Clone this repo: `git clone https://github.com/jicahoo/csi-unity.git`
* Create your branch: `git checkout -b <your_branch_name> --track origin/integration-hack-attack`
* Code changes on your branch. Add and commit your changes
* Push your changes to server to keep your work
    * `git push origin <your_branch_name>`
* build&test to make sure your code is OK before PR.
* Login GitHub and send PR to branch integration-hack-attack. 
    * Add reviewers if you want.
    * There will travis CI build check.
* Resolve conflicting if any.
* Merge it by yourself (after traivs build passed.)

## Dev env
* Any unix-like OS installed with Go can be used to develop this project.
* Below steps are just a possible option using docker to deveop this project.
1. Use git to clone code. You can use vim to edit these code.
    * `cd <your_code_dir>`
    * `git clone https://github.com/jicahoo/csi-unity.git`

2. Start golang docker and enter into it.
    * Get golang image: `docker pull golang`
    * Start a golang docker container in detached mode and mount your code path to it.
        * Create a container with name csi-unity`docker run -dti --name csi-unity -v <your_code_dir>/csi-unity:/go/src/github.com/jicahoo/csi-unity golang`
        * Note: The target path in container **MUST** be set as `/go/src/github.com/jicahoo/csi-unity`. Or you **CAN'T** start compile/run csi-unity successfully.
        * **Enter** into the container: `docker exec -ti csi-unity /bin/bash`

## How to build&run
* `cd $GOPATH/src/github.com/jicahoo/csi-unity`
* `go install`. This command will generate exe file $GOPATH/bin/csi-unity
* `export CSI_ENDPOINT=csi.sock`
* `$GOPATH/bin/csi-unity`. The command will start the csi-unity server.

## How to run unit test
* In command line, `go test ./...`

## How to test the function with CSI client: csc.
* Test tool: https://github.com/rexray/gocsi/tree/master/csc . csc is client of csi plugin.
* Install the test tool csc. 
    * `go get github.com/rexray/gocsi/csc`
    * You will find binary `csc` at $GOPATH/bin
* Prerequisite:
    * `export CSI_ENDPOINT=csi.sock`
* Commands:
    * `./csc controller list-volumes`

## More useful commands to start the csi-unity:
* `CSI_ENDPOINT=tcp://127.0.0.1:34555   X_CSI_REQ_LOGGING=true   X_CSI_REP_LOGGING=true   X_CSI_LOG_LEVEL=debug  go run main.go`

## Tools used by this project
* Go package dependency tool: https://github.com/golang/dep


## References about Go
### Go dep
* https://gist.github.com/subfuzion/12342599e26f5094e4e2d08e9d4ad50d
* https://blog.boatswain.io/post/manage-go-dependencies-using-dep/
* https://stackoverflow.com/questions/37237036/how-should-i-use-vendor-in-go-1-6

### Go tips
* go tags: https://flaviocopes.com/go-tags/ https://stackoverflow.com/questions/10858787/what-are-the-uses-for-tags-in-go
* make vs new: https://www.godesignpatterns.com/2014/04/new-vs-make.html https://stackoverflow.com/questions/9320862/why-would-i-make-or-new
