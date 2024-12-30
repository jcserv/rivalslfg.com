dev:
	go build ./cmd/rivalslfg/main.go && ./main

dev-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

dev-db-down:
	docker compose -p rivalslfg -f docker-compose.yml down -v

test-db:
	docker compose -p rivalslfg -f docker-compose.yml up --detach

test-db-down:
	docker compose -p rivalslfg -f docker-compose.test.yml down -v

clean:
	rm main
