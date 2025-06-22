FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o bin/app ./main.go

EXPOSE 8080

CMD ["./bin/app", "--mode", "server"]
