# Dockerfile for service1
FROM golang:latest

WORKDIR /app

RUN mkdir /usr/local/go/src/userTransactions

COPY ./users /usr/local/go/src/userTransactions/users

COPY ./logging /usr/local/go/src/userTransactions/logging

COPY ./cmd/usersd/ .
RUN ls -lah

RUN go build -o usersd main.go

CMD ["./usersd"]
