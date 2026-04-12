# Apache License 2.0
#
# Copyright 2006 Shane
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
SHELL := bash
ARTIFACT_NAME := golic
TESTPARALLELISM := 4
WORKING_DIR := $(shell pwd)
VERSION ?= v0.1.0
PACKAGE = github.com/Bugs5382/golic/internal/buildinfo
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifndef NO_COLOR
YELLOW=\033[0;33m
CYAN=\033[1;36m
RED=\033[31m
# no color
NC=\033[0m
endif

.PHONY: clean
clean::
	rm -rf $(WORKING_DIR)/bin

.PHONY: build
build:
	@mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X main.Version=$(VERSION)" \
		-o bin/$(ARTIFACT_NAME)-$(GOOS)-$(GOARCH)

.PHONY: test
test::
	go test -v -tags=all -parallel ${TESTPARALLELISM} -timeout 2h -covermode atomic -coverprofile=covprofile ./...

.PHONY: lint-init
lint-init:
	@echo -e "\n$(CYAN)Check for lint dependencies$(NC)"
	brew install golangci-lint
	brew install gitleaks
	brew install yamllint

.PHONY: lint
lint: test license
	@echo -e "\n$(YELLOW)Running the linters$(NC)"
	@echo -e "\n$(CYAN)golangci-lint$(NC)"
	goimports -w ./
	golangci-lint run
	@echo -e "\n$(CYAN)yamllint$(NC)"
	yamllint .
	@echo -e "\n$(CYAN)gitleaks$(NC)"
	gitleaks detect . --no-git --verbose --config=.gitleaks.toml


.PHONY: license
license: build
	./bin/golic-$(GOOS)-$(GOARCH) inject -c "2006 Shane" -t apache2

.PHONY: license-dry
license-dry: build
	./bin/golic-$(GOOS)-$(GOARCH) inject -c "2006 Shane" -t apache2 -d
