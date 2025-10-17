FROM golang:1.14-buster as build

ENV GO111MODULE="on"
WORKDIR /app
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download
COPY src .
RUN go build

# -----------

FROM debian:13.1

COPY --from=build /app/jen /jen
CMD ["bash"]
