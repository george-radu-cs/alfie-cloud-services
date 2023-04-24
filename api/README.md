# Alfie API

## Install dependencies

```bash
go mod tidy
go mod vendor
```

## Format code

```bash
go fmt ./...
```

## Generate protobuf files

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./app/protobuf/*.proto
```

```bash
cd app/protobuf && \
protoc --dart_out=grpc:../../app/dart_protobuf ./alfie_api.proto && \
cd ../../
```

## Build

```bash
go build -v -o ./build-destination ./...
```

> client on 
>> /Users/georgeradu/School/license/grpc-go-1.52.0/examples/helloworld
