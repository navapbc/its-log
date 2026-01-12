.PHONY: build clean generate serve docker test

build: clean generate
	go build -o its-log .

amd:
	GOOS=linux GOARCH=amd64 go build -o its-log . 

clean:
	rm -f its-log
	rm -rf internal/sqlite/models
	rm -f cmd/itslog/schema.sql

generate:
	cd internal/sqlite ; sqlc generate

serve: generate
	go run ./... serve

test:
	go test ./...
