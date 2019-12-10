export GO111MODULE=on

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
BUILDER_IMAGE=cvbarros/terraform-provider-teamcity-builder

default: test

build:
	GO111MODULE=on go build -o ./bin/terraform-provider-teamcity_${VERSION}

install: build
	cp ./bin/terraform-provider-teamcity_${VERSION} ~/.terraform.d/plugins/

clean:
	rm -rf ./bin

builder-action:
	docker run -e GITHUB_WORKSPACE='/github/workspace' -e GITHUB_REPOSITORY='terraform-provider-teamcity' -e GITHUB_REF='v0.0.1-alpha' --name terraform-provider-teamcity-builder $(BUILDER_IMAGE):latest

builder-image:
	docker build .github/builder --tag $(BUILDER_IMAGE)

clean_samples:
	find ./examples -name '*.tfstate' -delete
	find ./examples -name ".terraform" -type d -exec rm -rf "{}" \;

fmt_samples:
	terraform fmt -recursive examples/
