CURRENT_WORKING_DIR=$(shell pwd)

#------------------------------------------------------------------
# Project build information
#------------------------------------------------------------------
PROJNAME          		:= portcullis
VENDOR            		:= gavinmcnair
MAINTAINER        		:= Gavin McNair

GIT_REPO          		:= github.com/$(VENDOR)/$(PROJNAME)
GIT_SHA           		:= $(shell git rev-parse --verify HEAD)
BUILD_DATE        		:= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

DOCKER_USERNAME     	:= "gavinmcnair"
DOCKER_PASSWORD     	?="unknown"

# Construct docker image name.
IMAGE             		:= $(PROJNAME):$(VERSION)

#------------------------------------------------------------------
# Go configuration
#------------------------------------------------------------------
GOCMD             		:= go
GOFMT             		:= gofmt
BIN               		:= bin

#------------------------------------------------------------------
# Build targets
#------------------------------------------------------------------

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOFMT) -w main.go
	$(GOCMD) fmt *.go

.PHONY: vet
vet: ## Run go vet against code
	$(GOCMD) vet main.go
	$(GOCMD) vet .

.PHONY: test
test: fmt vet ## Run tests
	$(GOCMD) test . -coverprofile cover.out

.PHONY: portcullis
portcullis: fmt vet
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)
	env GOOS=darwin GOARCH=amd64 $(GOCMD) build -o $(BIN)/darwin/$(PROJNAME) $(GIT_REPO)

.PHONY: portcullis-linux
portcullis-linux: fmt vet
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)

#------------------------------------------------------------------
# CI targets
#------------------------------------------------------------------

.PHONY: build
build:
	docker build \
    --build-arg git_repository=`git config --get remote.origin.url` \
    --build-arg git_branch=`git rev-parse --abbrev-ref HEAD` \
    --build-arg git_commit=`git rev-parse HEAD` \
    --build-arg built_on=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    -t $(IMAGE) .

.PHONY: login
login:
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)
