# Используем образ Go на базе Alpine для сборки
FROM golang:1.22-alpine as builder

# Устанавливаем необходимые пакеты для поддержки CGo
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Устанавливаем переменную среды CGO_ENABLED в 1
ENV CGO_ENABLED=1

# Собираем приложение с поддержкой CGo
RUN go build -o main ./cmd/app

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем необходимые пакеты для запуска приложения с поддержкой CGo
RUN apk add --no-cache sqlite-libs

# Создаем директорию uploads
RUN mkdir -p /uploads

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

COPY --from=builder /app/.env /.env

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]