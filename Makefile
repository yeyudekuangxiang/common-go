SHELL := /bin/bash

test:
	cd ./tests && go test ./...
docker:
	cd ./scripts && sh docker-run.sh
lint:
	gofmt -w -s -l .