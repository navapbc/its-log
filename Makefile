.PHONY: amd native

native:
	cd api ; make native

e2e:
	cd api ; make e2e

k6-arm:
	cd ops/containers ; make k6-arm

k6-amd:
	cd ops/containers ; make k6-amd

itslog-arm:
	cd ops/containers ; make itslog-arm

itslog-arm-build:
	cd ops/containers ; make itslog-arm-build

itslog-amd:
	cmd ops/containers ; make itslog-amd

generate:
	cd api ; make generate