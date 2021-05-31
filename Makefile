
.PHONY: all build
build:
	go build -o ./bin/aliasd ./cmd/aliasd

all: build
