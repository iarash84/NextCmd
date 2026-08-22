GO ?= go
TARGET_DIR := target
HOST_OS := $(shell $(GO) env GOHOSTOS)
HOST_ARCH := $(shell $(GO) env GOHOSTARCH)
export GOCACHE := $(abspath $(TARGET_DIR)/.go-cache)

ifeq ($(HOST_OS),windows)
HOST_EXECUTABLE := NextCmd.exe
ROOT_EXECUTABLE := NextCmd.exe
define remove_artifacts
	@if exist "$(TARGET_DIR)" rmdir /s /q "$(TARGET_DIR)"
	@if exist "$(ROOT_EXECUTABLE)" del /q "$(ROOT_EXECUTABLE)"
endef
define copy_to_root
	@copy /Y "$(TARGET_DIR)\$(HOST_EXECUTABLE)" "$(ROOT_EXECUTABLE)" >NUL
endef
else
HOST_EXECUTABLE := NextCmd
ROOT_EXECUTABLE := NextCmd
define remove_artifacts
	@rm -rf "$(TARGET_DIR)"
	@rm -f "$(ROOT_EXECUTABLE)"
endef
define copy_to_root
	@cp "$(TARGET_DIR)/$(HOST_EXECUTABLE)" "$(ROOT_EXECUTABLE)"
endef
endif

NATIVE_OUTPUT := $(TARGET_DIR)/$(HOST_EXECUTABLE)
GO_SOURCES := main.go go.mod \
	$(wildcard sdk/*.go) \
	$(wildcard cmd/*/*.go) \
	$(wildcard internal/*/*.go) \
	$(wildcard plugins/*/*.go) \
	$(wildcard plugins/*/*/*.go) \
	$(wildcard examples/*/*.go)

.DEFAULT_GOAL := build

.PHONY: help build clean test test-race vet format run build-root build-all \
	build-windows-amd64 build-windows-arm64 \
	build-linux-amd64 build-linux-arm64 \
	build-darwin-amd64 build-darwin-arm64

help:
	@echo NextCmd Make targets:
	@echo   make help                  Show this command reference.
	@echo   make build                 Build for the current host into target/.
	@echo   make clean                 Remove target/ and the root executable.
	@echo   make test                  Run all unit tests.
	@echo   make test-race             Run race detector tests and requires CGO.
	@echo   make vet                   Run Go static analysis.
	@echo   make format                Format all Go packages.
	@echo   make run                   Build and run the current host executable.
	@echo   make build-root            Build for the host and copy the executable to the root.
	@echo   make build-all             Build all supported OS and architecture combinations.
	@echo Individual cross-build targets:
	@echo   make build-windows-amd64   Build Windows amd64.
	@echo   make build-windows-arm64   Build Windows arm64.
	@echo   make build-linux-amd64     Build Linux amd64.
	@echo   make build-linux-arm64     Build Linux arm64.
	@echo   make build-darwin-amd64    Build macOS amd64.
	@echo   make build-darwin-arm64    Build macOS arm64.

build: $(NATIVE_OUTPUT)

$(TARGET_DIR):
	@mkdir "$(TARGET_DIR)"

$(NATIVE_OUTPUT): $(GO_SOURCES) | $(TARGET_DIR)
	$(GO) build -trimpath -o "$@" .

clean:
	$(remove_artifacts)

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

vet:
	$(GO) vet ./...

format:
	$(GO) fmt ./...

run: build
	"$(NATIVE_OUTPUT)"

# Build only for the current host, then copy the result to the repository root.
build-root: build
	$(copy_to_root)

build-windows-amd64: export GOOS := windows
build-windows-amd64: export GOARCH := amd64
build-windows-amd64: export CGO_ENABLED := 0
build-windows-amd64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-windows-amd64.exe" .

build-windows-arm64: export GOOS := windows
build-windows-arm64: export GOARCH := arm64
build-windows-arm64: export CGO_ENABLED := 0
build-windows-arm64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-windows-arm64.exe" .

build-linux-amd64: export GOOS := linux
build-linux-amd64: export GOARCH := amd64
build-linux-amd64: export CGO_ENABLED := 0
build-linux-amd64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-linux-amd64" .

build-linux-arm64: export GOOS := linux
build-linux-arm64: export GOARCH := arm64
build-linux-arm64: export CGO_ENABLED := 0
build-linux-arm64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-linux-arm64" .

build-darwin-amd64: export GOOS := darwin
build-darwin-amd64: export GOARCH := amd64
build-darwin-amd64: export CGO_ENABLED := 0
build-darwin-amd64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-darwin-amd64" .

build-darwin-arm64: export GOOS := darwin
build-darwin-arm64: export GOARCH := arm64
build-darwin-arm64: export CGO_ENABLED := 0
build-darwin-arm64: | $(TARGET_DIR)
	$(GO) build -trimpath -o "$(TARGET_DIR)/NextCmd-darwin-arm64" .

build-all: build-windows-amd64 build-windows-arm64 \
	build-linux-amd64 build-linux-arm64 \
	build-darwin-amd64 build-darwin-arm64
