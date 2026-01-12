.PHONY: all
all: build

build:
	go build -o ./cmd/shortener ./cmd/shortener

test:
	go clean -testcache
	go test -count 1 -v ./...