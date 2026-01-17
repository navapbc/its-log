.PHONY: build clean generate serve docker prod stress k6 test

generate:
	cd api ; make generate

container-itslog:
	cd containers/itslog ; \
	docker build \
		--platform "linux/amd64" \
		-t itslog:latest \
		-f Dockerfile.itslog ../../itslog

container-jupyterlite:
	cd containers/jupyterlite ; \
	docker build \
		--platform "linux/amd64" \
		-t jupyterlite:latest \
		-f Dockerfile.jupyterlite .


e2e: amd
	@echo "e2e - root"
	cd containers/e2e ; make e2e

amd:
	cd api ; make amd

native:
	cd api ; make native

run: native
	./local.bash

up: 
	cd api ; make up

test:
	go test ./...
