.PHONY: k6 e2e itslog arm amd

amd:
	cd ../../api ; make linux-amd
arm:
	cd ../../api ; make linux-arm

k6:
	docker compose \
	-f docker-compose-k6.yaml \
	up \
	--abort-on-container-exit

e2e:
	docker compose \
	-f docker-compose-e2e.yaml \
	up \
	--abort-on-container-exit

k6-arm: arm k6

k6-amd: amd k6

itslog:
	docker compose \
	-f docker-compose.yaml \
	up

itslog-arm: arm itslog
