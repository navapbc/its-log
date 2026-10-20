.PHONY: build test stress
CONTAINER_CMD ?= $(shell command -v podman 2> /dev/null || echo docker)

prep-local:
	mkdir -p ~/.its-log/data

build: prep-local
	cd api ; go build -o its-log .

container-api: prep-local
	$(CONTAINER_CMD) build -t itslog-api:latest -f ops/containers/itslog/Dockerfile --load .

container-e2e: prep-local
	$(CONTAINER_CMD) build -t itslog-e2e:latest -f ops/containers/e2e/Dockerfile --load .

container-k6: prep-local
	$(CONTAINER_CMD) build -t itslog-k6 -f ops/containers/k6/Dockerfile --load .

e2e: container-api container-e2e
	$(CONTAINER_CMD) compose \
	-f ops/containers/docker-compose-e2e.yaml \
	up \
	--no-attach itslog \
	--abort-on-container-exit

e2e-noisy: container-api container-e2e
	$(CONTAINER_CMD) compose \
	-f ops/containers/docker-compose-e2e.yaml \
	up \
	--abort-on-container-exit

k6: container-api container-k6
	$(CONTAINER_CMD) compose \
	-f ops/containers/docker-compose-k6.yaml \
	up \
	--no-attach itslog \
	--abort-on-container-exit

run: container-api
	$(CONTAINER_CMD) compose \
	-f ops/containers/docker-compose.yaml \
	up \
	--abort-on-container-exit
