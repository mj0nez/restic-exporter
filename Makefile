MAKEFLAGS 		+= --warn-undefined-variables
SHELL 			:= bash
.SHELLFLAGS 	:= -eu -o pipefail -c
.DEFAULT_GOAL	:= help

# version:=go build -ldflags "-X main.version=`git tag --sort=-version:refname | head -n 1`
export APP_NAME		:= restic-exporter
export VERSION 		:= $(shell git tag --sort=-version:refname | head -n 1)
export REVISION 	:= $(shell git rev-parse HEAD)
export BUILD_DIR	:= ./bin/
# export STORE_DIR	:= ./.tmp/store

ifeq ($(CI),true)
	export GOTESTSUM_FORMAT := testname
endif

# check with:
# $ go build
# $ go tool nm ./tools.exe | grep version
# check correct path
#
# NOTE: the linker flags s and w remove debug symbols !
build:
	@go build \
		-v \
		-ldflags "\
			-s -w \
			-X 'github.com/mj0nez/restic-exporter/internal/info.Version=${VERSION}' \
			-X 'github.com/mj0nez/restic-exporter/internal/info.Revision=${REVISION}' \
			" \
		-o ${BUILD_DIR}
	@printf "Size:\t" && du -h $(BUILD_DIR)$(APP_NAME)
.PHONY: build

tests:
	@gotestsum
.PHONY: tests

bump-deps:
	@go get -u && go mod tidy
.PHONY: bump-deps

tools:
	@echo "===> Installing tools..."
	go install gotest.tools/gotestsum@v1.12.2
.PHONY: tools

tools-check:
	@echo "===> Checking installed tool versions..."
	go install github.com/Gelio/go-global-update@v0.2.5
	go-global-update --dry-run
.PHONY: tools-check

# runs the linter but excludes the dev directory

GOLANGCI_DIRS := ./internal ./cmd ./

lint:
	golangci-lint run ${GOLANGCI_DIRS}
.PHONY: lint

format:
	golangci-lint fmt ${GOLANGCI_DIRS}
	uvx pre-commit run --all-files
.PHONY: format
