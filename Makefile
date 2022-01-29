.PHONY: build
build:
	go build -o user-balance -v ./cmd/apiserver

run:
	./user-balance

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build