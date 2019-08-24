# Copyright 2019 The GoRE.tk Authors. All rights reserved.
# Use of this source code is governed by the license that
# can be found in the LICENSE file.

APP = libgore

SHELL = /bin/bash
DIR = $(shell pwd)
GO = go
UID=$(shell id -u)
GID=$(shell id -g)
DOCKER_FOLDER=docker
CONTAINER_NAME=gorebuild
VERSION=$(shell git describe --tags 2> /dev/null || git log --pretty=format:'%h' -n 1)

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
MAKE_COLOR=\033[33;01m%-20s\033[0m

APP_FILES=$(APP).so $(APP).dll $(APP).h
PACKAGE=$(APP)-$(VERSION)
LINUX_BUILD_FOLDER=build/linux
LINUX_ARCHIVE=$(PACKAGE)-linux-amd64.tar.gz
WINDOWS_ARCHIVE=$(APP)-$(VERSION)-windows.zip
WINDOWS_BUILD_FOLDER=build/windows
TAR_ARGS=cfz
RELEASE_FILES=LICENSE README.md

ARCH=GOARCH=amd64
CGO=CGO_ENABLED=1
BUILD_OPTS=-ldflags="-s -w" -buildmode=c-shared
WINDOWS_GO_ENV=GOOS=windows $(ARCH) $(CGO) CC=x86_64-w64-mingw32-gcc
LINUX_GO_ENV=GOOS=linux $(ARCH) $(CGO)

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo -e "$(OK_COLOR)==== $(APP) [$(VERSION)] ====$(NO_COLOR)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(MAKE_COLOR) : %s\n", $$1, $$2}'

.PHONY: windows
windows: ## Make binary for Windows
	@echo -e "$(OK_COLOR)[$(APP)] Build for Windows$(NO_COLOR)"
	@$(WINDOWS_GO_ENV) $(GO) build -o $(APP).dll $(BUILD_OPTS) .

.PHONY: linux
linux: ## Make binary for linux
	@echo -e "$(OK_COLOR)[$(APP)] Build for Linux$(NO_COLOR)"
	@$(LINUX_GO_ENV) $(GO) build -o $(APP).so $(BUILD_OPTS) .

.PHONY: build
build: ## Make binary
	@echo -e "$(OK_COLOR)[$(APP)] Build$(NO_COLOR)"
	@$(GO) build -o $(APP) $(BUILD_OPTS) .

.PHONY: clean
clean: ## Remove build artifacts
	@echo -e "$(OK_COLOR)[$(APP)] Clean$(NO_COLOR)"
	@rm -fr $(APP_FILES) build 2> /dev/null

.PHONY: docker_container
docker_container: ## Build build-container
	@echo -e "$(OK_COLOR)[$(APP)] Build docker container$(NO_COLOR)"
	@docker build -t $(CONTAINER_NAME):latest $(DOCKER_FOLDER)/

$(APP_FILES):
	@echo -e "$(OK_COLOR)[$(APP)] Build using docker container$(NO_COLOR)"
	@docker run -it --rm -u $(UID):$(GID) -v $(DIR):/go/libgore $(CONTAINER_NAME)
	@cat structs.h >> libgore.h

.PHONY: release

$(LINUX_ARCHIVE): $(APP).so $(APP).h
	@mkdir -p $(LINUX_BUILD_FOLDER)/$(PACKAGE)
	@cp $(RELEASE_FILES) $(APP).so $(APP).h $(LINUX_BUILD_FOLDER)/$(PACKAGE)/.
	@tar $(TAR_ARGS) $(LINUX_ARCHIVE) -C $(LINUX_BUILD_FOLDER) $(PACKAGE)

$(WINDOWS_ARCHIVE): $(APP).dll $(APP).h
	@mkdir -p $(WINDOWS_BUILD_FOLDER)/$(PACKAGE)
	@cp $(RELEASE_FILES) $(APP).dll $(APP).h $(WINDOWS_BUILD_FOLDER)/$(PACKAGE)/.
	@cd $(WINDOWS_BUILD_FOLDER) && zip -r $(DIR)/$(WINDOWS_ARCHIVE) $(PACKAGE) > /dev/null

release: $(LINUX_ARCHIVE) $(WINDOWS_ARCHIVE) ## Make release archives

docker_build: $(APP_FILES) ## Build using docker container

