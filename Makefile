.DEFAULT_GOAL := list

CONTAINER_REGISTRY ?= docker.io
CONTAINER_IMAGE ?= st3v/plotq
VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -o ./plotq ./cmd

run: build
	./plotq

test:
	go test -v --race ./... -count=1

clean:
	rm -f ./plotq

container-build:
	docker build -t ${CONTAINER_REGISTRY}/${CONTAINER_IMAGE}:${VERSION} .

container-push: container-build
	docker push ${CONTAINER_REGISTRY}/${CONTAINER_IMAGE}:${VERSION}

container-publish: container-build
	docker tag ${CONTAINER_REGISTRY}/${CONTAINER_IMAGE}:${VERSION} ${CONTAINER_REGISTRY}/${CONTAINER_IMAGE}:latest
	docker push ${CONTAINER_REGISTRY}/${CONTAINER_IMAGE}:latest

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
