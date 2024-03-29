PROJECT=gateway_kit


SVR_IMAGE=budda-gateway-kit
SVR_BIN=gateway-kit
MODE=debug

GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GO_VERSION ?= $(shell go version)
BUILD_TIME ?= $(shell date '+%Y-%m-%d__%H:%M:%S%p')
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
TAG := $(VERSION)__$(OS)_$(ARCH)
# Use linker flags to provide version/build settings
LDFLAGS="-X '$(PROJECT)/config.GoVersion=$(GO_VERSION)' -X $(PROJECT)/config.GitVersion=$(GIT_VERSION) -X $(PROJECT)/config.BuildTime=$(BUILD_TIME)"

.PHONY: all clean build mod docker-build

all: build

build:
	@cd ./src && go mod tidy
	-cd ./src && swag init -g ./cmd/server/main.go
	cd ./src/cmd/server && go build -ldflags ${LDFLAGS} -race -o ${SVR_BIN} main.go
	mv ./src/cmd/server/${SVR_BIN} ./deploy

mod:
	@cd ./src && go mod tidy


docker-build:
	@docker build -t ${SVR_IMAGE} --build-arg BINARY_NAME=${SVR_BIN} --build-arg MODE=${MODE}
