include .env
export 

.PHONY: codegen dev dev-db dev-db-down dev-db-migrate test-db test-db-down test-db-migrate clean

codegen:
	sqlc generate

dev:
	go build ./cmd/rivalslfg/main.go && ./main

dev-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

dev-db-down:
	docker compose -p rivalslfg -f docker-compose.yml down -v

migrate:
	migrate -database "$(DATABASE_URL)?sslmode=disable" -path ./db/migrations up

test-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

test-db-down:
	docker compose -p rivalslfg -f docker-compose.test.yml down -v

migrate-test:
	migrate -database "$(TEST_DATABASE_URL)?sslmode=disable" -path ./db/migrations up

clean:
	rm main
