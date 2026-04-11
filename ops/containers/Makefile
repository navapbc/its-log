.PHONY: e2e amd arm k6 kd-arm k6-amd itslog itslog-arm itslog-amd 

e2e:
	docker compose \
	-f docker-compose-e2e.yaml \
	up \
	--abort-on-container-exit

e2e-build:
	docker compose \
	-f docker-compose-e2e.yaml \
	up \
	--build \
	--abort-on-container-exit

amd:
	cd ../../api ; make linux-amd
arm:
	cd ../../api ; make linux-arm

native:
	cd ../../api ; make native

k6:
	docker compose \
	-f docker-compose-k6.yaml \
	up \
	--abort-on-container-exit

k6-arm: arm k6

k6-amd: amd k6

itslog:
	docker compose \
	-f docker-compose.yaml \
	up

itslog-build:
	docker compose \
	-f docker-compose.yaml \
	up --build


itslog-arm: arm itslog

itslog-arm-build: arm itslog-build

itslog-amd: amd itslog

itslog-amd-build: amd itslog-build