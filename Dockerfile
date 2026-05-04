FROM golang:1.24 AS builder

# Задаёт рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum отдельно (кеширование зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем 
# (CGO_ENABLED=0) полностью статический,не зависит от Linux-библиотек, работает в любом контейнере
# (RUN go build -o app) команда выполняется внутри /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/app


# 2. Финальный образ (минимальный)
FROM alpine:3.19

WORKDIR /app

# сертификаты (важно для HTTP)
RUN apk add --no-cache ca-certificates

# создаём пользователя
RUN adduser -D appuser

# Копируем только бинарник из builder
COPY --from=builder /app/app .

USER appuser

# Запуск
CMD ["./app"]