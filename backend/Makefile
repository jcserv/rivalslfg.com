include .env
export 

.PHONY: codegen dev clean dev-db dev-db-down dev-db-migrate reset test-db test-db-down test-db-migrate reset-test 

codegen:
	sqlc generate

dev:
	go build ./cmd/rivalslfg/main.go && ./main

clean:
	rm main

dev-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

dev-db-down:
	docker compose -p rivalslfg -f docker-compose.yml down -v

migrate:
	migrate -database "$(DATABASE_URL)?sslmode=disable" -path ./db/migrations up

reset:
	make dev-db-down && make dev-db

test-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

test-db-down:
	docker compose -p rivalslfg -f docker-compose.test.yml down -v

migrate-test:
	migrate -database "$(TEST_DATABASE_URL)?sslmode=disable" -path ./db/migrations up

migrate-prod:
	flyctl ssh console -C "migrate -database \"$DATABASE_URL\" -path /app/db/migrations up"

reset-test:
	make test-db-down && make test-db && make migrate-test
