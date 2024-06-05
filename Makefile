.PHONY: all
all: test build

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build
