.PHONY: build clean

build:
	go build -o ./bin/client ./cmd/client/
	go build -o ./bin/server ./cmd/server/

install:
	go install ./...

test:
	go test ./...

clean:
	rm -rf bin/*
