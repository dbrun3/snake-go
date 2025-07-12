.PHONY: build clean run run-client run-server deploy

ifneq (,$(wildcard .env))
  include .env
  export
endif

# Default value if not set in .env or shell
SERVER_URL ?= localhost:8080

build: 	go build -o bin/app ./main.go

clean:
	rm -rf bin/*

run:
	go run main.go --mode host

run-client:
	go run main.go --mode client --addr $(SERVER_URL)

run-server:
	go run main.go --mode server

OLD_VERSION := $(shell git describe --tags --abbrev=0)
NEW_VERSION := $(shell svu next)

release:
	@if [ "$$(uname)" = "Darwin" ]; then \
		sed -i "" "s/$(OLD_VERSION)/$(NEW_VERSION)/g" README.md; \
	else \
		sed -i "s/$(OLD_VERSION)/$(NEW_VERSION)/g" README.md; \
	fi; \
	git commit -m "chore: update readme release version"; \
	git tag $(svu next); \

