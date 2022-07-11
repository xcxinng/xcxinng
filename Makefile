# supported OS and ARCH, visit:
# https://go.dev/doc/install/source#environment
GOOS=darwin
GOARCH=arm64

GO=go
APP=go
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS="-X main.CommitId=$(COMMIT) -X main.Built=$(DATE) -X main.AppName=$(APP)"
GOBUILD=GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags $(LDFLAGS)

# setup go env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export GO111MODULE := on
export GOOS := $(GOOS)
export GOARCH := $(GOARCH)

consumer: APP=consumer
consumer: projects_in_action/self_monitoring/consumer.go
	$(GOBUILD) -v -o $(APP) projects_in_action/self_monitoring/consumer.go


producer: APP=producer
producer: projects_in_action/self_monitoring/producer.go
	$(GOBUILD) -v -o $(APP) projects_in_action/self_monitoring/producer.go
