.PHONY: all build consumer producer


all: build

build: consumer producer

consumer:
	go build -race -o bin/consumer -v consumer/main.go
producer:
	go build -race -o bin/producer -v producer/main.go