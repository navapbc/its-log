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


serve: generate
	cp config.yaml cmd
	go run ./... serve

docker:
	cd containers ; \
	docker build \
		--platform "linux/amd64" \
		-t its-log:latest \
		-f Dockerfile ..
	
prod:
	docker build -t its-log:latest -f Dockerfile --target prod ..

stress:
	cd k6 ; k6 run put.js
	ls -alh data/

k6:
	cd k6 ; k6 run put.js
	ls -alh data/

test:
	go test ./...