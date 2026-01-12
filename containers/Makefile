jsonnet:
	cd e2e ; make jsonnet

up: jsonnet
	docker compose up

e2e: jsonnet
	docker compose up e2e

test: jsonnet
	docker compose --profile testing up