# Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

PKG_DIR                     := $(REPO_ROOT)/pkg
TOOLS_DIR                   := $(REPO_ROOT)/tools


ifneq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	EFFECTIVE_VERSION := $(EFFECTIVE_VERSION)-dirty
endif

.DEFAULT_GOAL := all
all: watcher
#########################################
# Tools                                 #
#########################################

include $(REPO_ROOT)/hack/tools.mk
export PATH := $(abspath $(TOOLS_DIR)):$(PATH)


#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

.PHONY: docker-images
docker-images:
	@BUILD_ARCH=$(BUILD_ARCH) \
		$(REPO_ROOT)/hack/docker-image-build.sh "watcher" \
		$(IMAGE_REPOSITORY) $(IMAGE_TAG)

.PHONY: docker-push
docker-push:
	@$(REPO_ROOT)/hack/docker-image-push.sh "watcher" \
	$(IMAGE_REPOSITORY) $(IMAGE_TAG)

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

.PHONY: watcher
watcher:
	@echo "building $@ for $(BUILD_PLATFORM)/$(BUILD_ARCH)"
	@go mod tidy
	@go mod download
	@GOOS=$(BUILD_PLATFORM) \
		GOARCH=$(BUILD_ARCH) \
		CGO_ENABLED=0 \
		GO111MODULE=on \
		go build \
		-o $(REPO_ROOT)/build/watcher \
		-ldflags="$(LD_FLAGS)" \
		$(REPO_ROOT)/cmd/watcher

.PHONY: verify
verify: check test

.PHONY: verify-extended
verify-extended: check test-cov test-clean

.PHONY: clean
clean: test-clean
	@rm -f $(REPO_ROOT)/build/watcher

.PHONY: check
check: format $(GO_LINT)
	 @$(GO_LINT) run --config=$(REPO_ROOT)/.golangci.yaml --timeout 10m $(REPO_ROOT)/cmd/... $(REPO_ROOT)/pkg/...
	 @go vet $(REPO_ROOT)/cmd/... $(REPO_ROOT)/pkg/...

.PHONY: format
format:
	@gofmt -l -w $(REPO_ROOT)/cmd $(REPO_ROOT)/pkg

.PHONY: goimports
goimports: goimports_tool goimports-reviser_tool

.PHONY: goimports_tool
goimports_tool: $(GOIMPORTS)
	@for dir in $(SRC_DIRS); do \
		$(GOIMPORTS) -w $$dir/; \
	done

.PHONY: goimports-reviser_tool
goimports-reviser_tool: $(GOIMPORTS_REVISER)
	@for dir in $(SRC_DIRS); do \
		GOIMPORTS_REVISER_OPTIONS="-imports-order std,project,general,company" \
		$(GOIMPORTS_REVISER) -recursive $$dir/; \
	done

.PHONY: test
test:
	@go test $(REPO_ROOT)/cmd/parameters/... $(REPO_ROOT)/pkg/...

.PHONY: test-cov
test-cov:
	@$(REPO_ROOT)/hack/test-cover.sh ./cmd/... ./pkg/...

.PHONY: test-clean
test-clean:
	@rm -f $(REPO_ROOT)/test.*

.PHONY: add-license-headers
add-license-headers: $(GO_ADD_LICENSE)
	@./hack/add-license-header.sh