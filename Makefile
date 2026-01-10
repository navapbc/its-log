.PHONY: build clean generate serve docker prod stress k6 test

build: clean generate
	go build -o its-log .

clean:
	rm -f its-log
	rm -rf internal/sqlite/models
	rm -f cmd/itslog/config.yaml
	rm -f cmd/itslog/schema.sql

generate:
	cd internal/sqlite ; sqlc generate

config:
	mkdir -p ~/.itslog
	cp config.yaml ~/.itslog

serve: generate config
# 	go run ./... serve
	cd containers ; docker compose up

docker:
	cd containers ; \
	docker build \
		--platform "linux/amd64" \
		-t itslog:latest \
		-f Dockerfile ..
	
prod:
	docker build -t itslog:latest -f Dockerfile --target prod ..

stress:
	cd k6 ; k6 run put.js
	ls -alh data/

k6:
	cd k6 ; k6 run put.js
	ls -alh data/

test:
	go test ./...

etl: generate config
	cd pipeline ; make etl