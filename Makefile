METIS_VERSION := 0.0.1

build:
	go build -o bin/metis ./cmd/metis

proto:
	cd api && $(MAKE)

lint:
	golangci-lint run ./...

test:
	go test -v ./...

docker:
	docker build -t reg.navercorp.com/maas/metis-server:$(METIS_VERSION) -t reg.navercorp.com/maas/metis-server:latest .

.PHONY: build proto lint test docker
