jsonnet:
	jsonnet run.jsonnet > run.json
	jsonnet load.jsonnet > load.json
	jsonnet combine.jsonnet > combine.json

docker: jsonnet
	docker build \
		--platform "linux/amd64" \
		-t itslog-e2e:latest \
		-f Dockerfile.e2e .

e2e: jsonnet
	@echo E2E
	make jsonnet
	rm -f ../../its-log/data/itslog/*.sqlite 
	ITSLOG_DAYS=5 ITSLOG_EVENTS=1000 docker compose up e2e 
	ITSLOG_DAYS=1 ITSLOG_ACTIONS='combine' docker compose up e2e