export GO111MODULE = on

GIT_REF := $(shell git describe --always)
VERSION ?= $(GIT_REF)

ECHO_REPO := rerorero/echo
ECHO_IMAGE := $(ECHO_REPO):$(VERSION)

CONSUL_ENVYO_REPO := rerorero/consul-envoy
CONSUL_ENVOY_IMAGE := $(CONSUL_ENVYO_REPO):$(VERSION)

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: dependency
dependency:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go mod vendor

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:. ./proto/*.proto

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin/echo ./cmd/echo

.PHONY: container
container:
	docker build -t $(ECHO_IMAGE) -f docker/echo.build.Dockerfile .
	docker tag $(ECHO_IMAGE) $(ECHO_REPO):latest

.PHONY: dockerhub
dockerhub:
	docker login -u "$$DOCKER_USERNAME" -p "$$DOCKER_PASSWORD";
	docker push $(ECHO_IMAGE)
	docker push $(ECHO_REPO):latest

.PHONY: consul-envoy
consul-envoy:
	docker build -t $(CONSUL_ENVOY_IMAGE) -f docker/consul-envoy.Dockerfile .
	docker tag $(CONSUL_ENVOY_IMAGE) $(CONSUL_ENVYO_REPO):latest
	docker login -u "$$DOCKER_USERNAME" -p "$$DOCKER_PASSWORD";
	docker push $(CONSUL_ENVOY_IMAGE)
	docker push $(CONSUL_ENVYO_REPO):latest
