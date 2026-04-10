SHELL := bash
ARTIFACT_NAME := golic
TESTPARALLELISM := 4
WORKING_DIR := $(shell pwd)
GOLIC_VERSION  ?= v0.7.2

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
	go build -o bin/golic main.go

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
license:
	@echo -e "\n$(YELLOW)Injecting the license$(NC)"
	./bin/golic inject verbose -c "2022 Absa Group Limited" -t apache2
