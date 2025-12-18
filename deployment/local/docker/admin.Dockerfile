# ===== Этап сборки =====
FROM golang:1.24-alpine3.22 AS builder

# Устанавливаем необходимые утилиты
RUN apk --no-cache add ca-certificates git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и скачиваем зависимости
COPY services/admin/go.mod services/admin/go.sum ./
RUN go mod download

# Копируем исходники
COPY services/admin .

# Статическая сборка
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main ./cmd/main.go

# ===== Финальный минимальный образ =====
FROM alpine:3.22

# Устанавливаем только сертификаты
RUN apk --no-cache add ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.edge.kernel.org/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache tzdata


WORKDIR /app

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем бинарник
CMD ["./main"]
