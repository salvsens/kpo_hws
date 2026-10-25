# Единая точка входа для всех команд проекта.

.DEFAULT_GOAL := help

BIN := bin/vyshkat
GOLANGCI_LINT ?= golangci-lint

## help: показать это сообщение
.PHONY: help
help:
	@echo "Команды:"
	@sed -n 's/^## /  /p' $(MAKEFILE_LIST)

## run: запустить приложение
.PHONY: run
run:
	go run ./cmd/vyshkat

## build: собрать бинарник в bin/vyshkat
.PHONY: build
build:
	go build -o $(BIN) ./cmd/vyshkat

## test: запустить тесты
.PHONY: test
test:
	go test ./...

## test-race: тесты с детектором гонок
.PHONY: test-race
test-race:
	go test -race ./...

## cover: тесты с отчётом о покрытии
.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## fmt: отформатировать код (gofumpt + goimports из конфига)
.PHONY: fmt
fmt:
	$(GOLANGCI_LINT) fmt

## lint: прогнать линтеры (.golangci.yml)
.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

## vet: go vet
.PHONY: vet
vet:
	go vet ./...

## tidy: причесать go.mod/go.sum
.PHONY: tidy
tidy:
	go mod tidy

## check: полная проверка перед сдачей: fmt + vet + lint + test
.PHONY: check
check: fmt vet lint test
