.PHONY: build-amd generate up test e2e

build-amd:
	GOOS=linux GOARCH=amd64 go build -o its-log . 

generate:
	cd internal/sqlite ; sqlc generate

up: build-amd
	cd ../containers ; docker compose up

test:
	cd itslog ; go test ./...
