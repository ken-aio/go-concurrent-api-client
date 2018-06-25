# ref: https://github.com/dtan4/s3url/blob/master/Makefile
NAME := go-concurrent-api-client
VERSION := 1.0.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep
	dep ensure -v

.PHONY: update-deps
update-deps: dep
	dep ensure -update -v

.PHONY: run
run:
	go run main.go

.PHONY: run-mock
run-mock: dep
ifeq ($(shell command -v gin 2> /dev/null),)
	go get -u github.com/codegangsta/gin
endif
	cd mock && dep ensure -v
	cd mock && gin -a 9998 -p 9999 -i --notifications run api-server.go

.PHONY: build
build: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -v -o artifacts/$(NAME)
