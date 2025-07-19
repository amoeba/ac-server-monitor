.PHONY: all build test clean seed

all: build

build:
	go build -o monitor app.go
	go build -o seed cmd/seed/main.go

seed: build
	go build -o seed cmd/seed/main.go

test:
	go test ./...

clean:
	rm -f monitor seed *.db

.DEFAULT_GOAL := build