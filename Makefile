.PHONY: build clean

build: 	go build -o bin/client-app ./cmd/client

clean:
	rm -rf bin/*