.PHONY: build
build:
	docker-compose build user-balance

run:
	docker-compose up user-balance

stop:
	docker-compose down

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:0000@localhost:5436/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:0000@localhost:5436/postgres?sslmode=disable' down

.PHONY: test
test:
	go test -v -race -timeout 30s ./...