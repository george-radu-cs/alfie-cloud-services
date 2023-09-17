# Alfie API

Main project api. Contains the api for auth and media mangement. Written in Go. Uses gRPC for communication.

## Environment variables

Copy the example file:

```bash
cp .env.example .env
```

Or if using docker compose:

```bash
cp .env.example .env.docker
```

## Install dependencies

```bash
go mod tidy
go mod vendor
```

## Format code

```bash
# less strict
# go fmt ./...
# stricter (recommended)
gofumpt -l -w .
```

## Generate protobuf files

### For the go service

```bash
protoc --go_out=./internal/pkg/protobuf/ --go-grpc_out=./internal/pkg/protobuf/ ./api/protobuf/*.proto
```

### For the dart client

```bash
protoc --dart_out=grpc:./internal/pkg/dart_protobuf ./api/protobuf/*.proto
```
