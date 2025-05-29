.PHONY: build build-client build-server clean

build: build-client build-server

build-client:
	go build -o bin/client-app ./cmd/client

build-server:
	go build -o bin/server-app ./cmd/server

clean:
	rm -rf bin/*