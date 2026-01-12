.PHONY: build clean generate serve docker prod stress k6 test

e2e: generate
	cd containers ; docker compose --profile testing up

container-itslog:
	cd containers/itslog ; \
	docker build \
		--platform "linux/amd64" \
		-t itslog:latest \
		-f Dockerfile.itslog ../../itslog

container-e2e:
	cd containers/e2e ; \
	make docker

amd:
	cd itslog ; make amd

up: amd
	cd containers ; make up

up-test: amd
	cd containers ; make test

test:
	go test ./...
