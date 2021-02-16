build:
	go build -o bin/metis ./cmd/metis

test:
	go test ./...

.PHONY: build test
