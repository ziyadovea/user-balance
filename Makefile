.PHONY: build
build:
	go build -o user-balance -v ./cmd/apiserver

run:
	./user-balance

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:0000@localhost:5436/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:0000@localhost:5436/postgres?sslmode=disable' down

.DEFAULT_GOAL := build