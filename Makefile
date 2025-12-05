SQLC=sqlc
clean:
	rm -f itslog
	rm -rf internal/sqlite/models
	rm -f cmd/itslog/config.yaml
	rm -f cmd/itslog/schema.sql

generate:
	cd internal/sqlite ; ${SQLC} generate

build: clean generate
	go build -o itslog ./cmd/itslog/

run: generate
	cp config.yaml cmd/itslog/
	cd cmd/itslog ; go run ./...

docker:
	docker build -t itslog:latest -f Dockerfile .