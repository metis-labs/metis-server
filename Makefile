build:
	go build -o bin/metis ./cmd/metis

proto:
	cd api && $(MAKE)

test:
	go test ./...

.PHONY: build proto test
