# Подключение переменных из .env
-include .env
export

# Флаги для линкера (оптимизация размера)
LDFLAGS := -s -w

# Путь до main
MAIN := ./cmd/main.go

# Выходной бинарник
BIN := ./tmp/main

# Путь к Swagger точке входа
SWAG_ENTRY := ./cmd/main.go

# Путь к миграциям
MIGRATIONS_DIR := ./migrations

# Команда по умолчанию
.PHONY: all
all: mod sqlc lint swag build

# Обновление зависимостей
.PHONY: mod
mod:
	@echo "tidy dependencies..."
	go mod tidy

# Генерация SQLC-кода
.PHONY: sqlc
sqlc:
	@echo "generating SQLC code..."
	sqlc generate

# Проверка кода линтером
.PHONY: lint
lint:
	@echo "linting..."
	golangci-lint run ./...

# Сборка приложения
.PHONY: build
build:
	@echo "building..."
	go build -ldflags="${LDFLAGS}" -buildvcs=false -o $(BIN) $(MAIN)

# Запуск dev-приложения с reflex (live reload)
.PHONY: run
run:
	@echo "starting reflex live reload..."
	reflex -r '\.go$$' -s -- sh -c 'go build -o $(BIN) $(MAIN) && $(BIN)'

# Генерация Swagger-документации
.PHONY: swag
swag:
	@echo "generating Swagger docs..."
	swag init --parseDependency --parseInternal --parseDepth 5 -g $(SWAG_ENTRY)

# Очистка билда и временных файлов
.PHONY: clean
clean:
	@echo "cleaning..."
	rm -rf tmp

# Запуск всех тестов
.PHONY: test
test:
	@echo "running tests..."
	go test -v ./...

# Docker Compose запуск
.PHONY: up
up:
	@echo "starting docker-compose..."
	docker-compose -f docker-compose.yml up -d --build

# Docker Compose остановка
.PHONY: down
down:
	@echo "stopping docker-compose..."
	docker-compose down
