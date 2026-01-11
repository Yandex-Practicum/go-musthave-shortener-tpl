.PHONY: all
all: build

build:
	go build -o ./cmd/shortener ./cmd/shortener

my-test:
	./shortenertest -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener