BASEDIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NETWORK = docker_default

build:
	docker run --rm -v "$(BASEDIR)":/usr/src/mqtt_sync -w /usr/src/mqtt_sync golang:alpine \
		apk add --update --no-cache git openssh && \
		go get -d ./... && \
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./docker/mqtt_sync mqtt_sync.go
	docker build -t nonsenz/mqtt_sync ./docker/
	rm -f ./docker/mqtt_sync

clean:
	rm -f ./docker/mqtt_sync