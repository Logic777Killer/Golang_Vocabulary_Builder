# Этап 1: Сборка
FROM golang:1.25.8-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем основной сервер
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Собираем админ-утилиту
RUN CGO_ENABLED=0 GOOS=linux go build -o admin ./cmd/admin

# Этап 2: Минимальный образ
FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/admin .
COPY --from=builder /app/web ./web

USER appuser

EXPOSE 8080

# По умолчанию запускаем сервер
CMD ["./server"]