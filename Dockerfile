# =========================
# 1. Builder stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# зависимости системы (для CGO / сертификатов)
RUN apk add --no-cache git ca-certificates

# копируем модули сначала (для кеша)
COPY go.mod go.sum ./
RUN go mod download

# копируем весь проект
COPY . .

# сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/app

# =========================
# 2. Runtime stage
# =========================
FROM alpine:latest

WORKDIR /app

# сертификаты (TLS, HTTPS, etc.)
RUN apk add --no-cache ca-certificates

# копируем бинарник из builder
COPY --from=builder /app/auth-service .

# копируем migrations (если нужны внутри контейнера)
COPY --from=builder /app/migrations ./migrations

# порт приложения
EXPOSE 8080

# запуск
CMD ["./auth-service"]