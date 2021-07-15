PKG=github.com/iychoi/stock-svc
VERSION=v0.1.0
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO111MODULE=on
GOPROXY=direct
GOPATH=$(shell go env GOPATH)
SERVER_ADDR=18.144.20.11

.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=1 GOOS=linux go build ./cmd/main.go

.PHONY: deploy
deploy: build
	scp -r main exec resources ${SERVER_ADDR}:~/

.PHONY: deploy_resources
deploy_resources:
	scp -r resources ${SERVER_ADDR}:~/
