# Используем официальный образ Go для сборки приложения
FROM golang:1.22-bullseye AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Собираем бинарный файл
RUN go build -o main .

# Создаем финальный минимальный образ для запуска
FROM debian:bullseye-slim

# Устанавливаем рабочую директорию внутри финального контейнера
WORKDIR /app

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/main .

# Указываем переменные окружения по умолчанию (их можно переопределить при запуске контейнера)
ENV DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=admin \
    DB_NAME=ozzcafe

# Открываем порт для приложения (например, 8080)
EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["./main"]
