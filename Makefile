METIS_VERSION := 0.0.4

GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date "+%Y-%m-%d")
GO_PROJECT = oss.navercorp.com/metis/metis-server
GO_SRC := $(shell find . -path ./vendor -prune -o -type f -name '*.go' -print)

# inject the version number into the golang version package using the -X linker flag
GO_LDFLAGS ?=
GO_LDFLAGS += -X ${GO_PROJECT}/internal/version.GitCommit=${GIT_COMMIT}
GO_LDFLAGS += -X ${GO_PROJECT}/internal/version.Version=${METIS_VERSION}
GO_LDFLAGS += -X ${GO_PROJECT}/internal/version.BuildDate=${BUILD_DATE}

build:
	go build -o bin/metis -ldflags "${GO_LDFLAGS}" ./cmd/metis

proto:
	cd api && $(MAKE)

fmt:
	gofmt -s -w $(GO_SRC)

lint:
	golangci-lint run ./...

test:
	go test -v ./...

docker:
	docker build -t reg.navercorp.com/metis/metis-server:$(METIS_VERSION) -t reg.navercorp.com/metis/metis-server:latest .

.PHONY: build proto lint test docker
