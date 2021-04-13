METIS_VERSION := 0.0.1

GO_SRC := $(shell find . -path ./vendor -prune -o -type f -name '*.go' -print)

build:
	go build -o bin/metis ./cmd/metis

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
