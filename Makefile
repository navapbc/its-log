.PHONY: amd native

amd:
	cd api ; make amd

native:
	cd api ; make native

run:
	cd api ; make run

swagger:
	cd api ; swag init -o ./docs

generate:
	cd api ; make generate