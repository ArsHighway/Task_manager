# Task Manager API (Go + Gin + PostgreSQL)

Минимальный REST API для управления пользователями и задачами.

## Что умеет проект

- CRUD для задач (`/tasks`)
- CRUD для пользователей (`/users`)
- Получение задач конкретного пользователя (`GET /users/:id/tasks`)
- Фильтрация/сортировка/пагинация списка задач через query-параметры
- Хранение данных в PostgreSQL
- Автоприменение миграций в Docker через `migrate/migrate`

## Технологии

- Go 1.25
- [gin](https://github.com/gin-gonic/gin)
- [pgx/v5](https://github.com/jackc/pgx)
- [squirrel](https://github.com/Masterminds/squirrel)
- PostgreSQL 15
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- Docker Compose

## Быстрый старт (Docker)

```bash
docker compose up --build
```

После запуска:
- API: `localhost:8080`
- PostgreSQL: `localhost:5432`

Переменная окружения приложения:
- `DATABASE_URL`

По умолчанию в compose используется:
- `postgres://postgres:postgres@postgres:5432/app?sslmode=disable`

## Запуск локально (без Docker)

1) Подними PostgreSQL локально и создай БД `app`  
2) Примени миграции (например через `migrate`)  
3) Запусти приложение:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/app?sslmode=disable"
go run ./cmd/app
```

Сервер слушает `:8080`.

## REST API

### Tasks

- `POST /tasks/` — создать задачу
- `GET /tasks/` — получить список задач
- `GET /tasks/:id` — получить задачу по id
- `PUT /tasks/:id` — полное обновление задачи
- `PATCH /tasks/:id` — частичное обновление задачи
- `DELETE /tasks/:id` — удалить задачу

### Users

- `POST /users/` — создать пользователя
- `GET /users/:id` — получить пользователя по id
- `PATCH /users/:id` — частично обновить пользователя
- `DELETE /users/:id` — удалить пользователя
- `GET /users/:id/tasks` — получить задачи пользователя

## Пример JSON-моделей

### User

```json
{
  "id": 1,
  "name": "Ars",
  "email": "ars@example.com",
  "created_at": "2026-05-05T15:00:00Z"
}
```

### Task

```json
{
  "id": 10,
  "title": "Сделать README",
  "description": "Описать API и запуск",
  "status": "pending",
  "user_id": 1,
  "created_at": "2026-05-05T15:05:00Z"
}
```

## Фильтры и пагинация для `GET /tasks/`

Поддерживаются query-параметры:

- `status` — фильтр по статусу
- `user_id` — фильтр по пользователю
- `created_at` — фильтр по времени создания
- `sort_by` — поле сортировки (`created_at`, `title`, `status`)
- `order_by` — направление (`ASC` или `DESC`)
- `limit` — размер страницы (1..100, по умолчанию `10`)
- `page` — номер страницы (по умолчанию `1`)

Пример:

```text
GET /tasks/?status=pending&user_id=1&sort_by=created_at&order_by=DESC&limit=10&page=1
```

## Миграции

Миграции лежат в папке `migrations/`:

- создание таблицы `users`
- создание таблицы `tasks` (FK на `users`, `ON DELETE CASCADE`)

## Структура проекта

```text
cmd/app/main.go
internal/config/
internal/handlers/
internal/repository/
internal/routers/
internal/service/
internal/models/
internal/errs/
migrations/
docker-compose.yaml
```

## Полезные команды

```bash
go test ./...
go build ./...
docker compose up --build
docker compose down
```

## Ограничения текущей версии

- Нет аутентификации/авторизации
- Нет отдельного `health` endpoint
- Нет OpenAPI/Swagger-спецификации
- Ошибки возвращаются как текстовые сообщения, а не единый JSON-формат
