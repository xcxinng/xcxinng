# By default, GO use go configured by path env.
GO=go

# By default, all binary and archive files will be stored here.
InstallPath=$(shell pwd)

# Project dir, do not change this path(unless you know what you're doing).
projectPath=$(shell pwd)

# Compilation related information.
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS="-X main.CommitId=$(COMMIT) -X main.Built=$(DATE)"
GOBUILD=$(GO) build -ldflags $(LDFLAGS) -v

# Setup go env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export GO111MODULE := on

consumer:
	$(GOBUILD) -o $(InstallPath)/consumer $(ProjectPath)/action_abnormal_event/consumer/consumer.go

producer:
	$(GOBUILD) -o $(InstallPath)/producer ./action_abnormal_event/producer/producer.go

archive_consumer: consumer
	tar -czvf $(InstallPath)/consumer.tar.gz $(InstallPath)/consumer $(ProjectPath)/action_abnormal_event/consumer/consumer.service

archive_producer: producer
	tar -czvf $(InstallPath)/producer.tar.gz $(InstallPath)/producer $(ProjectPath)/action_abnormal_event/consumer/producer.service

clean:
	$(GO) clean
	rm -f $(InstallPath)/consumer $(InstallPath)/producer consumer.tar.gz producer.tar.gz
