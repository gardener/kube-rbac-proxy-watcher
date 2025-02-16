# SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

NAME                        := kube-rbac-proxy-watcher
REGISTRY                    ?= europe-docker.pkg.dev/gardener-project/snapshots/gardener/extensions
IMAGE_REPOSITORY            := $(REGISTRY)/kube-rbac-proxy-watcher
REPO_ROOT                   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN                         := $(REPO_ROOT)/bin
VERSION                     := $(shell cat "$(REPO_ROOT)/VERSION")
EFFECTIVE_VERSION           := $(VERSION)-$(shell git rev-parse HEAD)
SRC_DIRS                    := $(shell go list -f '{{.Dir}}' $(REPO_ROOT)/...)
LD_FLAGS                    := $(shell $(REPO_ROOT)/hack/get-build-ld-flags.sh)
BUILD_PLATFORM              ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
BUILD_ARCH                  ?= $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
IMAGE_TAG                   := $(VERSION)

ifneq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	EFFECTIVE_VERSION := $(EFFECTIVE_VERSION)-dirty
endif

.DEFAULT_GOAL := all
all: watcher

#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

docker-images:
	@BUILD_ARCH=$(BUILD_ARCH) \
		$(REPO_ROOT)/hack/docker-image-build.sh "watcher" \
		$(IMAGE_REPOSITORY) $(IMAGE_TAG)

docker-push:
	@$(REPO_ROOT)/hack/docker-image-push.sh "watcher" \
	$(IMAGE_REPOSITORY) $(IMAGE_TAG)

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

tidy:
	@go mod tidy
	@go mod download

watcher: tidy
	@echo "building $@ for $(BUILD_PLATFORM)/$(BUILD_ARCH)"
	@GOOS=$(BUILD_PLATFORM) \
		GOARCH=$(BUILD_ARCH) \
		CGO_ENABLED=0 \
		GO111MODULE=on \
		go build \
		-o $(REPO_ROOT)/build/watcher \
		-ldflags="$(LD_FLAGS)" \
		$(REPO_ROOT)/cmd/watcher

verify: check test sast

verify-extended: check test-cov sast-report

clean: test-clean
	@rm -f $(REPO_ROOT)/build/watcher

check: tidy format
	 @go tool golangci-lint run \
	 	--config=$(REPO_ROOT)/.golangci.yaml \
		--timeout 10m \
		$(SRC_DIRS)
	 @go vet \
	 	$(SRC_DIRS)

format:
	@gofmt -l -w $(SRC_DIRS)

goimports: goimports_tool goimports-reviser_tool

goimports_tool: tidy
	@for dir in $(SRC_DIRS); do \
		go tool goimports -w $$dir/; \
	done

goimports-reviser_tool:
	@for dir in $(SRC_DIRS); do \
		GOIMPORTS_REVISER_OPTIONS="-imports-order std,project,general,company" \
		go tool goimports-reviser -recursive $$dir/; \
	done

test:
	@go tool gotest.tools/gotestsum $(SRC_DIRS)

test-cov:
	@$(REPO_ROOT)/hack/test-cover.sh $(SRC_DIRS)

test-clean:
	@rm -f $(REPO_ROOT)/test.*

add-license-headers: tidy
	@./hack/add-license-header.sh

sast: tidy
	@$(REPO_ROOT)/hack/sast.sh

sast-report: tidy
	@$(REPO_ROOT)/hack/sast.sh --gosec-report true

.PHONY: docker-images docker-push watcher verify verify-extended clean check format goimports goimports_tool goimports-reviser_tool test test-cov test-clean add-license-headers sast sast-report
