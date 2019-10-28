export GO111MODULE = on

GIT_REF := $(shell git describe --always)
VERSION ?= $(GIT_REF)

ECHO_REPO := rerorero/echo
ECHO_IMAGE := $(ECHO_REPO):$(VERSION)

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: build
	CGO_ENABLED=0 go build -o bin/echo ./cmd/echo

.PHONY: container
container:
	docker build -t $(ECHO_IMAGE) -f docker/echo.build.Dockerfile .
	docker tag $(ECHO_IMAGE) $(ECHO_REPO):latest

.PHONY: dockerhub
dockerhub:
	docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD";
	docker push $(ECHO_IMAGE)
	docker push $(ECHO_REPO):latest
