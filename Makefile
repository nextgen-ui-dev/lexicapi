.PHONY: build

install:
	go mod tidy

run:
	go run main.go

dev:
	air

build:
	go build -o bin/main main.go

run-build:
	./bin/main

up:
	docker compose up -d

down:
	docker compose down

migration:
	migrate create -seq -ext sql -dir db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate:
	go run db/migrations/migrate.go $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	go run db/migrations/migrate.go -action down $(filter-out $@,$(MAKECMDGOALS))
