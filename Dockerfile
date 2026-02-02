FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект и собираем его
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o forum ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/forum .

EXPOSE 8080

CMD ["./forum"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/forum .

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./forum"]