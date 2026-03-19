include .env
export

MIGRATIONS_DIR=migrations

.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

.PHONY: migrate-new
migrate-new:
ifndef name
	$(error Please provide migration name: make migrate-new name=<name>)
endif
	goose -dir $(MIGRATIONS_DIR) create $(name) sql