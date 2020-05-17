DOCKER_REGISTRY ?= docker-registry.dbld.tech
DOCKER_USER ?= build
IMAGE    ?= mindworks-software/iotwifi
NAME     ?= txwifi
VERSION  ?= 1.0.6

all: build push

dev: dev_build dev_run

build:
	docker build -t $(IMAGE):arm32v6-$(VERSION) .
	docker tag $(IMAGE):arm32v6-$(VERSION) $(IMAGE):latest

push:
	docker login -u $(DOCKER_USER) $(DOCKER_REGISTRY)
	docker tag $(IMAGE):latest $(DOCKER_REGISTRY)/$(IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE):latest
	docker tag $(IMAGE):arm32v6-$(VERSION) $(DOCKER_REGISTRY)/$(IMAGE):arm32v6-$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(IMAGE):arm32v6-$(VERSION)

dev_build:
	docker build -t $(IMAGE) ./dev/

dev_run:
	sudo docker run --rm -it --privileged --network=host \
                   -v $(CURDIR):/go/src/github.com/txn2/txwifi \
                   -w /go/src/github.com/txn2/txwifi \
                   --name=$(NAME) $(IMAGE):latest


