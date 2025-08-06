# syntax=docker/dockerfile:1.4
FROM golang:1.24.4 AS dev

WORKDIR /app

# Добавляем go bin в PATH, иначе reflex не найдется
ENV PATH="$PATH:/go/bin"

# Установка reflex
RUN go install github.com/cespare/reflex@latest

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Устанавливаем swag, если нужен
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Указываем точку входа
CMD ["reflex", "-r", "\\.go$$", "-s", "--", "sh", "-c", "go build -o tmp/main ./cmd/main.go && ./tmp/main"]
