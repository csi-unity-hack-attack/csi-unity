REGISTRY_NAME=ciqihuo
#REV=$(shell git describe --long --tags --dirty)
REV=1024
IMAGE_NAME=csi-unity
IMAGE_VERSION=$(REV)
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: all driver clean driver-container

all: driver

test:
	go test github.com/jicahoo/csi-unity/... -cover
	go vet github.com/jicahoo/csi-unity/...
driver:
	if [ ! -d ./vendor ]; then dep ensure -vendor-only; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-X github.com/jicahoo/csi-unity/main.vendorVersion=$(REV) -extldflags "-static"' -o _output/csi-unity .
driver-container: driver
	docker build -t $(IMAGE_TAG) .
push: driver-container
	docker push $(IMAGE_TAG)
clean:
	go clean -r -x
	-rm -rf _output
