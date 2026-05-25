# Этап 1: Сборка
FROM golang:1.25.8-alpine AS builder

WORKDIR /app

# Копируем ВСЁ сразу (внешних зависимостей нет, кэшировать нечего)
COPY . .

# Собираем бинарник (просто go build, без лишних флагов)
RUN go build -o server ./cmd/server

# Этап 2: Минимальный образ для запуска
FROM alpine:3.19

# Создаём пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Копируем только готовый бинарник и фронтенд
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

# Запускаем от обычного пользователя
USER appuser

# Порт приложения
EXPOSE 8080

# Команда запуска
CMD ["./server"]