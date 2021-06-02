
.PHONY: all build
build:
	go build -o ./bin/aliasd ./cmd/aliasd

all: build

.PHONY: proto
proto: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/aliasd.proto

