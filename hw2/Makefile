# Команды проекта «ПОС».

.DEFAULT_GOAL := help

BIN := bin/pos
GOLANGCI_LINT ?= golangci-lint

## help: показать команды
.PHONY: help
help:
	@sed -n 's/^## /  /p' $(MAKEFILE_LIST)

## run: запустить приложение
.PHONY: run
run:
	go run ./cmd/pos

## build: собрать бинарник в bin/pos
.PHONY: build
build:
	go build -o $(BIN) ./cmd/pos

## test: запустить тесты
.PHONY: test
test:
	go test ./...

## test-race: запустить тесты с детектором гонок
.PHONY: test-race
test-race:
	go test -race ./...

## cover: запустить тесты и показать покрытие
.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## fmt: отформатировать код
.PHONY: fmt
fmt:
	$(GOLANGCI_LINT) fmt

## lint: прогнать линтеры
.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

## lint-config: проверить конфигурацию линтера
.PHONY: lint-config
lint-config:
	$(GOLANGCI_LINT) config verify

## vet: проверить код через go vet
.PHONY: vet
vet:
	go vet ./...

## tidy: синхронизировать go.mod и go.sum с импортами
.PHONY: tidy
tidy:
	go mod tidy

## check: отформатировать код и выполнить vet, lint и test
.PHONY: check
check: fmt
	$(MAKE) vet lint test
