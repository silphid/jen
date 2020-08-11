FROM golang:1.14-buster

ENV GO111MODULE=on
WORKDIR /app
COPY . .
ENV go build
