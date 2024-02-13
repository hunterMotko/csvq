# Change these variables as necessary.
MAIN := ./cmd/main.go
BINARY_NAME := csvq

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## build: build the application
.PHONY: build
build:
	go build -o=./${BINARY_NAME} ${MAIN}
