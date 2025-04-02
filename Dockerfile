# Используем официальный образ Go для сборки
FROM golang:1.24 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл без зависимостей от C (статическая линковка)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kanban ./cmd/server

# Используем минимальный контейнер
FROM alpine:latest

# Устанавливаем необходимые утилиты
RUN apk --no-cache add ca-certificates

# Задаем рабочую директорию
WORKDIR /root/

# Копируем бинарник из builder-контейнера
COPY --from=builder /app/kanban ./

# Даем файлу права на выполнение
RUN chmod +x ./kanban

# Открываем порт
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["./kanban"]