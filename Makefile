.PHONY: build clean install test

build:
	go build -o ./bin/local/client ./cmd/client/
	go build -o ./bin/local/server ./cmd/server/

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/linux/client ./cmd/client/
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/linux/server ./cmd/server/

install:
	go install ./...

test:
	go test ./...

clean: docker-stop
	rm -rf bin/*

.PHONY:docker-start
docker-start:
	docker-compose up

.PHONY:docker-stop
docker-stop:
	docker-compose down