# URL Shortener

Сервис для сокращения ссылок на Go. Поддерживает два типа хранилища:

- `memory` - хранение в памяти приложения;
- `postgres` - хранение в PostgreSQL.

Короткая ссылка имеет длину 10 символов и состоит из символов `a-z`, `A-Z`, `0-9`, `_`.

## Требования

- Go 1.22+
- Docker и Docker Compose

## Запуск через Docker Compose

Основной способ запуска проекта:

```bash
docker compose up --build
```

Compose поднимает два контейнера:

- `db` - PostgreSQL на порту `5432`;
- `app` - HTTP API на порту `8080`.

Приложение запускается с переменными окружения:

```env
STORAGE=postgres
POSTGRES_DSN=postgres://user:pass@db:5432/shortener?sslmode=disable
PORT=8080
```

Остановить контейнеры:

```bash
docker compose down
```

Остановить контейнеры и удалить volume с данными PostgreSQL:

```bash
docker compose down -v
```

## Локальный запуск без Docker

По умолчанию сервис использует storage `memory` и порт `8080`:

```bash
go run ./cmd/app
```

Запуск с PostgreSQL:

```bash
STORAGE=postgres \
POSTGRES_DSN="postgres://user:pass@localhost:5432/shortener?sslmode=disable" \
PORT=8080 \
go run ./cmd/app
```

Перед локальным запуском с PostgreSQL база должна быть доступна. Можно поднять только базу через Docker Compose:

```bash
docker compose up db
```

Таблица `urls` создается приложением автоматически при старте.

## API

### Создать короткую ссылку

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://google.com"}'
```

Пример ответа:

```json
{
  "short": "abcDEF123_"
}
```

Для одного и того же оригинального URL сервис возвращает одну и ту же короткую ссылку.

### Перейти по короткой ссылке

```bash
curl -i http://localhost:8080/<short>
```

Например:

```bash
curl -i http://localhost:8080/abcDEF123_
```

Если короткая ссылка найдена, сервис вернет HTTP redirect `302 Found` на оригинальный URL.

## Тестирование

Запустить все тесты:

```bash
go test ./...
```

Запустить тесты сервиса:

```bash
go test -v ./cmd/internal/service
```

Если Go не может писать в системный build cache, используйте локальный кэш проекта:

```bash
mkdir -p .gocache
GOCACHE="$PWD/.gocache" go test -v ./...
```

## Структура проекта

```text
cmd/app/main.go                         # точка входа приложения
cmd/internal/config/config.go           # конфигурация из env
cmd/internal/handler/handlers.go        # HTTP handlers
cmd/internal/service/shortlink_gen.go   # бизнес-логика
cmd/internal/repository/                # memory и postgres repositories
cmd/internal/utils/generator.go         # генерация короткой ссылки
cmd/internal/service/shortlink_gen_test.go
Dockerfile
docker-compose.yml
```

## Переменные окружения

| Переменная | Значение по умолчанию | Описание |
| --- | --- | --- |
| `PORT` | `8080` | HTTP порт приложения |
| `STORAGE` | `memory` | Тип хранилища: `memory` или `postgres` |
| `POSTGRES_DSN` | пусто | DSN подключения к PostgreSQL |

`POSTGRES_DSN` обязателен, если `STORAGE=postgres`.
