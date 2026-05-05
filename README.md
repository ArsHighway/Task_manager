## Запуск

### Через Docker

```bash
docker pull arshighway/task-manager

docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://user:password@host:5432/dbname \
  arshighway/task-manager
```

### Локально

**Требования:** Go 1.24+, PostgreSQL

```bash
git clone https://github.com/ArsHighway/Task_manager.git
cd Task_manager
go mod download
go run main.go
```

> По умолчанию приложение подключается к `postgres://arsver@localhost:5432/arsver`.

## API

Сервер запускается на порту `:8080`.

### Задачи `/tasks`

| Метод    | Путь         | Описание                        |
|----------|--------------|----------------------------------|
| `POST`   | `/tasks/`    | Создать задачу                  |
| `GET`    | `/tasks/`    | Получить список задач (фильтры) |
| `GET`    | `/tasks/:id` | Получить задачу по ID           |
| `PUT`    | `/tasks/:id` | Полное обновление задачи        |
| `PATCH`  | `/tasks/:id` | Частичное обновление задачи     |
| `DELETE` | `/tasks/:id` | Удалить задачу                  |

### Пользователи `/users`

| Метод    | Путь                | Описание                          |
|----------|---------------------|-----------------------------------|
| `POST`   | `/users/`           | Создать пользователя              |
| `GET`    | `/users/:id`        | Получить пользователя по ID       |
| `GET`    | `/users/:id/tasks`  | Получить все задачи пользователя  |
| `PATCH`  | `/users/:id`        | Частичное обновление пользователя |
| `DELETE` | `/users/:id`        | Удалить пользователя              |

## Примеры запросов

### Создать пользователя
```bash
curl -X POST http://localhost:8080/users/ \
  -H "Content-Type: application/json" \
  -d '{"name": "Ivan", "email": "ivan@example.com"}'
```

### Создать задачу
```bash
curl -X POST http://localhost:8080/tasks/ \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко", "user_id": 1}'
```

### Частичное обновление задачи
```bash
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Новое название"}'
```

### Получить все задачи пользователя
```bash
curl http://localhost:8080/users/1/tasks
```

## Docker Hub

Образ доступен на Docker Hub: [`arshighway/task-manager`](https://hub.docker.com/repository/docker/arshighway/task-manager)

## Лицензия

MIT
