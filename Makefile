# SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

NAME                        := kube-rbac-proxy-watcher
REPO_ROOT                   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
REGISTRY                    ?= europe-docker.pkg.dev/gardener-project/snapshots/gardener/extensions
IMAGE_REPOSITORY            := $(REGISTRY)/kube-rbac-proxy-watcher
VERSION                     := $(shell cat "$(REPO_ROOT)/VERSION")
EFFECTIVE_VERSION           := $(VERSION)-$(shell git rev-parse --short HEAD)
SRC_DIRS                    := $(shell go list -f '{{.Dir}}' $(REPO_ROOT)/...)
BUILD_PLATFORM              ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
BUILD_ARCH                  ?= $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
LD_FLAGS                    ?= "-s -w $(shell $(REPO_ROOT)/hack/get-ldflags.sh)"

GCI_OPT                     ?= -s standard -s default -s "prefix($(shell go list -m))" --skip-generated

ifneq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	EFFECTIVE_VERSION       := $(EFFECTIVE_VERSION)-dirty
endif
IMAGE_TAG                   := $(EFFECTIVE_VERSION)

.DEFAULT_GOAL := all
all: check watcher

#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

docker-images:
	@docker build \
  		--build-arg LD_FLAGS=$(LD_FLAGS) \
  		--tag $(IMAGE_REPOSITORY):latest \
		--tag $(IMAGE_REPOSITORY):$(IMAGE_TAG) \
  		--platform linux/$(BUILD_ARCH) \
		-f Dockerfile --target watcher $(REPO_ROOT)

docker-push:
	@docker push $(IMAGE_REPOSITORY):latest
	@docker push $(IMAGE_REPOSITORY):$(IMAGE_TAG)

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: gci
gci: tidy
	@echo "Running gci..."
	@go tool gci write $(GCI_OPT) $(SRC_DIRS)

.PHONY: watcher
watcher: tidy
	@echo "building $@ for $(BUILD_PLATFORM)/$(BUILD_ARCH)"
	@GOOS=$(BUILD_PLATFORM) \
		GOARCH=$(BUILD_ARCH) \
		CGO_ENABLED=0 \
		go build \
			-o $(REPO_ROOT)/build/watcher \
			-ldflags=$(LD_FLAGS) \
			$(REPO_ROOT)/cmd/watcher

.PHONY: fmt
fmt: tidy
	@echo "Running $@..."
	@go tool golangci-lint fmt \
    	--config=$(REPO_ROOT)/.golangci.yaml \
    	$(SRC_DIRS)

.PHONY: check
check: tidy fmt gci lint

.PHONY: lint
lint: tidy
	@echo "Running $@..."
	 @go tool golangci-lint run \
	 	--config=$(REPO_ROOT)/.golangci.yaml \
		$(SRC_DIRS)

.PHONY: test
test:
	@echo "Running $@..."
	@go tool gotest.tools/gotestsum $(SRC_DIRS)

.PHONY: test-cov
test-cov:
	@echo "Running $@..."
	@$(REPO_ROOT)/hack/test-cover.sh $(SRC_DIRS)

.PHONY: verify
verify: check test sast

.PHONY: verify-extended
verify-extended: check test-cov sast-report

.PHONY: clean
clean:
	@echo "Running $@..."
	@rm -f $(REPO_ROOT)/build/watcher
	@rm -f $(REPO_ROOT)/*.sarif
	@rm -f $(REPO_ROOT)/test.coverprofile
	@rm -f $(REPO_ROOT)/test.coverage.html

sast: tidy
	@echo "Running $@..."
	@$(REPO_ROOT)/hack/sast.sh

sast-report: tidy
	@echo "Running $@..."
	@$(REPO_ROOT)/hack/sast.sh --gosec-report true

add-license-headers: tidy
	@echo "Running $@..."
	@$(REPO_ROOT)/hack/add-license-header.sh
