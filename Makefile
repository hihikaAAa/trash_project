.PHONY: build run test clean swagger docker-build check

# Имя бинарного файла
BINARY_NAME=trash_project

# Сборка сервиса
build:
	go build -o out/$(BINARY_NAME) ./cmd/server

# Запуск сервиса
run:
	go run ./cmd/server/main.go

# Запуск тестов
test:
	go test -v ./...

# Очистка
clean:
	rm -rf out/
	rm -rf docs/

# Генерация Swagger документации
swagger:
	swag init --dir cmd/server --pdl 3 --output docs

# Сборка Docker образа
docker-build:
	docker build -t $(BINARY_NAME):latest .

# Установка зависимостей
deps:
	go mod tidy
	go mod vendor

# Проверки линтеров + доп.генерация
check:
	gocritic check ./...
	fieldalignment -fix ./...
	golangci-lint run ./...
	go mod tidy
	go fix ./...
	swag init --dir cmd/server --pdl 3 --output docs
