FROM golang:latest

WORKDIR /usr/src/api

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading 
# them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy the source code
COPY . .
# copy the env file to the root of the project
COPY .env.docker .env

# install protoc compiler
RUN apt update && apt install -y protobuf-compiler
# install the go plugins for protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# generate protobuf files
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./app/protobuf/*.proto

# build the app
RUN go build -v -o /usr/local/bin ./...

CMD ["app"]