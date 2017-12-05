.PHONY: build build-linux clean install test

all: ensure install test

build: gobindata
	go build -o ./bin/local/client ./cmd/client/
	go get -u github.com/jteeuwen/go-bindata/
	go-bindata -o ./app/server/bindata.go -pkg server ./app/server/static/...
	go build -o ./bin/local/server ./cmd/server/

build-linux: gobindata
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/linux/client ./cmd/client/
	go-bindata -o ./app/server/bindata.go -pkg server ./app/server/static/...
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/linux/server ./cmd/server/

ensure:
	dep ensure

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

gobindata:
	go get -u github.com/jteeuwen/go-bindata/...