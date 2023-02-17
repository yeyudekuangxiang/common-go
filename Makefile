SHELL := /bin/bash

test:
	cd ./tests && go test ./...
docker:
	cd ./scripts && sh docker-run.sh
lint:
	gofmt -w -s -l .
vet:.
	go vet ./cmd/mp2c ./cmd/consumer ./cmd/command