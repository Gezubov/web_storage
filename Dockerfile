# Используем образ Go на базе Alpine для сборки
FROM golang:1.22-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Собираем приложение
RUN go build -o main ./cmd/app

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

COPY --from=builder /app/.env /.env

COPY --from=builder app/swagger.json /docs/swagger.json

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]