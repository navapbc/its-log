.PHONY: amd generate up test e2e

clean:
	rm -f its-log

linux-amd: clean generate
	GOOS=linux GOARCH=amd64 go build -o its-log . 

linux-arm: clean generate
	GOOS=linux GOARCH=arm64 go build -o its-log . 

native: clean generate
	go build -o its-log . 

run: native
	./local.bash

generate:
	cd internal/schema ; sqlc generate

up: build-amd
	cd ../containers ; docker compose up itslog lite

test:
	go test ./...
