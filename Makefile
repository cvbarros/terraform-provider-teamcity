.PHONY: build build-alpine clean test help default

export GO111MODULE=on
BIN_NAME=terraform-provider-teamcity

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
BUILDER_IMAGE=cvbarros/terraform-provider-teamcity-builder

default: test

help:
	@echo 'Management commands for sandbox:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'

	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/cvbarros/terraform-provider-teamcity/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/cvbarros/terraform-provider-teamcity/version.BuildDate=${BUILD_DATE}" -o ${BIN_NAME}

builder-action:
	docker run -e GITHUB_WORKSPACE='/github/workspace' -e GITHUB_REPOSITORY='terraform-provider-teamcity' -e GITHUB_REF='v0.0.1-alpha' --name terraform-provider-teamcity-builder $(BUILDER_IMAGE):latest

builder-image:
	docker build .github/builder --tag $(BUILDER_IMAGE)

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

