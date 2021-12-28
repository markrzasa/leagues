THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

compose-up:
	docker-compose -f $(THIS_DIR)/compose/docker-compose.yml up -d

build:
	go build -o out/ $(THIS_DIR)

dependencies:
	go get -u
	go mod tidy -compat=1.17

run:
	go run $(THIS_DIR)

test:
	go test $(THIS_DIR)
