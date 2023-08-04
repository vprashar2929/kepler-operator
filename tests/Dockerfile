FROM golang:1.20 AS builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
COPY ./tests ./tests
RUN go mod download

RUN go test -c ./tests/ -o mytest

FROM registry.access.redhat.com/ubi9-minimal:9.2-691

WORKDIR /workspace

COPY --from=builder /go/src/app/mytest .

# ENTRYPOINT [ "./mytest" ]