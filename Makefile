.PHONY: build clean run run-client run-server

SERVER_URL=localhost:8080

build: 	go build -o bin/app ./main.go

clean:
	rm -rf bin/*

run:
	go run main.go

run-client:
	go run main.go --mode client --addr $(SERVER_URL)

run-server:
	go run main.go --mode server