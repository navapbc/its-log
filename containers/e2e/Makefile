jsonnet:
	jsonnet run.jsonnet > run.json
	jsonnet load.jsonnet > load.json

docker: jsonnet
	docker build \
		--platform "linux/amd64" \
		-t itslog-e2e:latest \
		-f Dockerfile.e2e .
