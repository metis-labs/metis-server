build:
	go build -o bin/metis ./cmd/metis

proto:
	cd api && $(MAKE)

lint:
	 golangci-lint run ./...

test:
	go test -v ./...

.PHONY: build proto test
