.PHONY: build clean generate serve docker prod stress k6 test

generate:
	cd itslog ; make generate

container-itslog:
	cd containers/itslog ; \
	docker build \
		--platform "linux/amd64" \
		-t itslog:latest \
		-f Dockerfile.itslog ../../itslog

e2e:
	@echo "e2e - root"
	cd itslog ; make up & \
	cd ../containers/e2e ; make e2e


amd:
	cd itslog ; make build-amd

up:
	cd itslog ; make up

test:
	go test ./...
