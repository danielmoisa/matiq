.PHONY: build all test clean

all: build

run:
	go run ./cmd/auto-runner/main.go

run-web:
	cd web && npm run dev

build: build-http-server build-websocket-server build-http-server-internal

build-http-server:
	go build -o bin/auto-runner src/cmd/auto-runner/main.go

build-websocket-server:
	go build -o bin/auto-runner-websocket src/cmd/auto-runner-websocket/main.go

docker-compose:
	docker-compose up -d

swagger:
	swag init --parseDependency --parseInternal -g ./cmd/auto-runner/main.go -o docs


test:
	PROJECT_PWD=$(shell pwd) go test -race ./...

test-cover:
	go test -cover --count=1 ./...

cover-total:
	go test -cover --count=1 ./... -coverprofile cover.out
	go tool cover -func cover.out | grep total 

cov:
	PROJECT_PWD=$(shell pwd) go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

fmt:
	@gofmt -w $(shell find . -type f -name '*.go' -not -path './*_test.go')

fmt-check:
	@gofmt -l $(shell find . -type f -name '*.go' -not -path './*_test.go')

init-database:
	/bin/bash scripts/postgres-init.sh

clean:
	@ro -fR bin
