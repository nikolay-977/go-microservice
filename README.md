# Go Microservice - Управление пользователями

Микросервис для управления пользователями с поддержкой высоких нагрузок.

## Функциональность

- RESTful API для CRUD операций с пользователями
- Rate limiting (1000 запросов/секунду)
- Метрики Prometheus
- Асинхронное логирование
- Уведомления через worker pool
- Контейнеризация через Docker

## API Endpoints

### Пользователи
- `GET /api/users` - список всех пользователей
- `GET /api/users/{id}` - получить пользователя по ID
- `POST /api/users` - создать пользователя
- `PUT /api/users/{id}` - обновить пользователя
- `DELETE /api/users/{id}` - удалить пользователя

### Системные
- `GET /health` - проверка здоровья
- `GET /metrics` - метрики Prometheus

## Запуск

### Локально
```bash
make build
make run
```

```bash
make clean
```

### Запуск в Docker

```bash
make docker-run
```

```bash
make docker-stop
```

## Postman коллекции и переменные

[text](<postman/Go Microservice API.postman_collection.json>)
[text](<postman/Go Microservice Environment.postman_environment.json>)

## Нагрузочное тестирование
```bash
make load-test
```
