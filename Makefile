SQLC=sqlc

build: clean generate
	go build -o its-log .

clean:
	rm -f its-log
	rm -rf internal/sqlite/models
	rm -f cmd/itslog/config.yaml
	rm -f cmd/itslog/schema.sql

generate:
	cd internal/sqlite ; ${SQLC} generate


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
	cd k6 ; k6 run get.js
	ls -alh data/
	cd k6 ; k6 run put.js
	ls -alh data/