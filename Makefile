MIGRATIONS_DIR := ./db/migrations

DATABASE_URL := "postgres://postgres:postgres@localhost:5432/machine?sslmode=disable"

.PHONY: run
run:
	@go run ./cmd/server/main.go

.PHONY: migrate-new
migrate-new:
ifndef name
	$(error Você deve especificar o nome da migração. Exemplo: make migrate-new name=create_users_table)
endif
	@echo "==> Criando nova migração: $(name)"
	@touch $(MIGRATIONS_DIR)/`date +%Y%m%d%H%M%S`_$(name).up.sql
	@touch $(MIGRATIONS_DIR)/`date +%Y%m%d%H%M%S`_$(name).down.sql

.PHONY: migrate-up
migrate-up:
	@echo "==> Aplicando migrações..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up


.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  run: Executa o programa"