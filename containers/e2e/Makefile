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
	rm -f ../../data/itslog/*.sqlite 
	docker compose up itslog &
	ITSLOG_DAYS=14 ITSLOG_EVENTS=500 docker compose up e2e 
	ITSLOG_EVENTS=1 ITSLOG_DAYS=1 ITSLOG_ACTIONS='combine' docker compose up e2e
	docker stop containers-itslog-1