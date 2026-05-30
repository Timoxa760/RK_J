.PHONY: up down build test migrate migrate-down migrate-ch seed clean

POSTGRES_URL=postgres://postgres:postgres@localhost:5432/moneymind?sslmode=disable

up:
	docker compose up -d

down:
	docker compose down

build:
	go build -v -o ./bin/api-gateway.exe ./services/core-api/api-gateway/
	go build -v -o ./bin/user-service.exe ./services/core-api/user-service/
	go build -v -o ./bin/receipt-service.exe ./services/receipt-engine/receipt-service/
	go build -v -o ./bin/scraper-service.exe ./services/receipt-engine/scraper-service/
	go build -v -o ./bin/category-service.exe ./services/finance-core/category-service/
	go build -v -o ./bin/budget-service.exe ./services/finance-core/budget-service/
	go build -v -o ./bin/credit-service.exe ./services/finance-core/credit-service/
	go build -v -o ./bin/bank-service.exe ./services/finance-core/bank-service/
	go build -v -o ./bin/ai-processor.exe ./services/money-intelligence/ai-processor/
	go build -v -o ./bin/analytics-service.exe ./services/money-intelligence/analytics-service/
	go build -v -o ./bin/gamification.exe ./services/social-game/gamification/
	go build -v -o ./bin/social-service.exe ./services/social-game/social-service/
	go build -v -o ./bin/notification-service.exe ./services/reporting/notification-service/
	go build -v -o ./bin/reporting-service.exe ./services/reporting/reporting-service/
	@echo All 14 services built.

test:
	go test -v ./internal/... ./services/...

migrate:
	docker run --rm --network moneymind_network \
		-v $(CURDIR)/db/migrations/postgres:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_URL)" up

migrate-down:
	docker run --rm --network moneymind_network \
		-v $(CURDIR)/db/migrations/postgres:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_URL)" down 1

migrate-ch:
	docker run --rm --network moneymind_network \
		-v $(CURDIR)/db/migrations/clickhouse:/migrations \
		clickhouse/clickhouse-server:25.12 \
		clickhouse-client --host clickhouse --user clickhouse_user --password clickhouse_password \
		--queries-file /migrations/001_receipt_items.sql

seed:
	go run scripts/seed_data.go

clean:
	rm -f ./bin/*.exe
