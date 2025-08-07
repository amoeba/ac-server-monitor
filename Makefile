.PHONY: all build test clean seed

all: build

build:
	go build -ldflags "-X monitor/lib.GitHash=$$(git rev-parse --short HEAD)" -o monitor app.go
	go build -o seed cmd/seed/main.go

test:
	go test ./...

clean:
	rm -f monitor seed
