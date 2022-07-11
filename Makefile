GO=go
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS="-X main.CommitId=$(COMMIT) -X main.Built=$(DATE)"
GOBUILD=$(GO) build -ldflags $(LDFLAGS) -v

# setup go env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export GO111MODULE := on

consumer:
	$(GOBUILD) -o consumer ./action_abnormal_event/consumer/consumer.go

producer:
	$(GOBUILD) -o producer ./action_abnormal_event/producer/producer.go

clean:
	$(GO) clean
	rm -f consumer producer
