# Тестовое задание на позицию Junior Go Developer (Effective Mobile)

## Содержание

- [Задание](#задание)
  - [Требования](#требования)
    - [Общие требования](#общие-требования)
    - [Технические требования](#технические-требования)
- [Реализация](#реализация)
  - [Возможности API](#возможности-api)
  - [Технологии](#технологии)
  - [Запуск проекта](#запуск-проекта)
  - [Переменные окружения](#переменные-окружения)
  - [Swagger](#swagger)
  - [Тестирование](#тестирование)
  - [Структура проекта](#структура-проекта)

## Задание

Реализовать REST-сервис для агрегации данных об онлайн-подписках пользователей.

### Требования

#### Общие требования

Сервис должен:

- предоставлять CRUDL-операции над записями о подписках;
- считать суммарную стоимость подписок с фильтрами;
- работать с PostgreSQL;
- иметь документацию Swagger.
- покрыт логами

#### Технические требования

Сервис реализован на [Go](https://go.dev/) с PostgreSQL, миграциями и контейнеризацией через Docker Compose.

## Реализация

Проект представляет собой HTTP API для управления подписками, реализован с помощью фреймфорка [Gin](https://github.com/gin-gonic/gin).

### Возможности API

- `POST /subscriptions` - создать подписку
- `GET /subscriptions` - получить список подписок
- `GET /subscriptions/{id}` - получить подписку по ID
- `PATCH /subscriptions/{id}` - частично обновить подписку
- `DELETE /subscriptions/{id}` - удалить подписку
- `GET /subscriptions/total` - получить суммарную стоимость подписок

### Технологии

- Go `1.26.3`
- Gin
- PostgreSQL 17
- pgx
- golang-migrate
- zap logger
- swaggo (Swagger)
- Docker / Docker Compose

### Запуск проекта

1. Скопировать `.env.example` в `.env`:

```bash
cp .env.example .env
```

2. Заполнить `.env` и проверить порты:
- `HTTP_PORT` - порт API на хосте (например, `8080`)
- `PG_PORT` - порт PostgreSQL на хосте (например, `5432`)

3. Для запуска в Docker использовать:
- `PG_HOST=pg`
- `HTTP_HOST=0.0.0.0`

4. Запустить:

```bash
docker compose up --build
```

5. API будет доступен на порту из `HTTP_PORT`.

### Переменные окружения

Пример (`.env`):

```env
PG_USERNAME=admin
PG_PASSWORD=1234
PG_DATABASE=subscribe-manager
PG_HOST=pg
PG_PORT=5432

HTTP_HOST=0.0.0.0
HTTP_PORT=8080

LOGGER_LEVEL=dev
```

### Swagger

После запуска документация доступна по адресу:

- `http://localhost:8080/swagger/index.html`

Также реализован редирект:

- `GET /` -> `/swagger/index.html`

Если нужно перегенерировать swagger-спеку:

```bash
swag init -g cmd/main/main.go -d cmd/main,internal -o docs
```

```

### Тестирование

- Слои `internal/service` и `internal/handler` покрыты тестами на **100%**.
- Запуск всех тестов:

```bash
go test ./...
```

### Структура проекта

```text
cmd
└── main
    └── main.go                 # Точка входа

internal
├── app
│   └── app.go                  # Bootstrap приложения (config, db, migrations, http)
├── apperrors
│   └── err.go                  # Доменные ошибки приложения
├── config
│   ├── http.go                 # HTTP config
│   ├── load.go                 # Загрузка .env
│   ├── logger.go               # Конфиг логгера
│   └── pg.go                   # Конфиг PostgreSQL
├── db
│   ├── db.go                   # Обертка над pgxpool
│   └── interface.go
├── domain
│   └── subscription.go         # Доменные сущности
├── dto
│   └── dto.go                  # Request/Response модели
├── handler
│   ├── middleware
│   │   └── logger.go           # HTTP логирование + request_id
│   ├── subscription            # HTTP-ручки подписок + swagger аннотации
│   ├── handler.go
│   └── router.go               # Инициализация роутов
├── logger
│   └── logger.go
├── repository
│   ├── postgres                # Реализация репозитория на PostgreSQL
│   ├── mocks
│   ├── generate.go
│   └── repository.go
└── service
    ├── subscription            # Бизнес-логика
    ├── mocks
    ├── generate.go
    └── service.go

migrations
├── 001_create_subscriptions.up.sql
└── 001_create_subscriptions.down.sql

docs                            # Сгенерированная swagger документация
Dockerfile
docker-compose.yaml
.env.example
```
