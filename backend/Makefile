include .env
export 

.PHONY: codegen dev mocks clean dev-db dev-db-down dev-db-migrate reset test-db test-db-down test-db-migrate reset-test depgraph

codegen:
	sqlc generate

dev:
	go build ./cmd/rivalslfg/main.go && ./main

mocks:
	mockgen -package=mocks -source=internal/services/interfaces.go -destination=internal/test/mocks/interfaces.go

test:
	go test ./... 

clean:
	rm main

dev-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

dev-db-down:
	docker compose -p rivalslfg -f docker-compose.yml down -v

migrate:
	migrate -database "$(DATABASE_URL)?sslmode=disable" -path ./db/migrations up

# Usage: make makemigration name=your_migration_name
migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

reset:
	make dev-db-down && make dev-db

test-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

test-db-down:
	docker compose -p rivalslfg -f docker-compose.test.yml down -v

migrate-test:
	migrate -database "$(TEST_DATABASE_URL)?sslmode=disable" -path ./db/migrations up

reset-test:
	make test-db-down && make test-db && make migrate-test

depgraph:
	godepgraph github.com/jcserv/rivalslfg/internal | dot -Tpng -o godepgraph.png
