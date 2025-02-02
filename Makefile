.PHONY: run
run:
	@go run ./cmd/server/main.go


.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  run: Executa o programa"