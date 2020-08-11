FROM golang:1.14-buster as build

ENV GO111MODULE="on"
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src src
RUN go build

# -----------

FROM debian:buster

COPY --from=build /app/jen /jen
CMD ["bash"]
