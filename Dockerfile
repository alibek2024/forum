# Stage 1: Сборка
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Собираем приложение (CGO_ENABLED=0 для pure Go sqlite если нужно, но modernc работает и так)
RUN CGO_ENABLED=0 GOOS=linux go build -o forum ./cmd/app

# Stage 2: Финальный образ
FROM alpine:latest
# ca-certificates нужны, sqlite в самом образе не обязателен (у тебя драйвер внутри бинарника)
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 1. Копируем бинарник
COPY --from=builder /app/forum .

# 2. Копируем папку с миграциями (в ней уже лежит твой schema.sql)
COPY --from=builder /app/migrations ./migrations

# 3. Копируем фронтенд
COPY --from=builder /app/web ./web

# 4. Копируем конфиг (он есть у тебя на скрине)
COPY --from=builder /app/config.yaml .

EXPOSE 8081
CMD ["./forum"]