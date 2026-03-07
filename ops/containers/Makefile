.PHONY: k6

e2e:
	cd e2e ; make e2e

amd:
	cd ../../api ; make linux-amd
arm:
	cd ../../api ; make linux-arm

k6:
	docker compose \
	-f docker-compose-k6.yaml \
	up \
	--abort-on-container-exit


k6-arm: arm k6

k6-amd: amd k6
