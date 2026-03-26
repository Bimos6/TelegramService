# Telegram Service

Сервис для управления множественными соединениями с Telegram через gRPC API.

---

## Архитектура

<img width="968" height="603" alt="image" src="https://github.com/user-attachments/assets/7e3298b3-778e-42b1-a3ef-00ece72ceebe" />


### Ключевые архитектурные решения

| Решение | Описание |
|---------|----------|
| **Изоляция сессий** | Каждая сессия работает в отдельной горутине, ошибки в одной не влияют на другие |
| **In-memory хранилище** | Для тестирования, легко заменяется на Redis/PostgreSQL |
| **Graceful shutdown** | Корректное завершение всех соединений при остановке сервера |
| **Интерфейс логгера** | Бизнес-логика не зависит от конкретной реализации (logrus) |
| **Mock клиент** | Возможность тестирования без реальных ключей Telegram |

---

## 🚀 Быстрый старт

### Требования
- Go 1.26+
- grpcurl или Postman (для тестирования)

### Установка и запуск

```bash
# 1. Клонирование репозитория
git clone https://github.com/Bimos6/TelegramService.git
cd TelegramService 

# 2. Копирование .env
cp .env.example .env

# 3. Установка зависимостей
go mod download

# 4. Генерация protobuf кода
cd .\internal\app\grpc\proto\
protoc --go_out=. --go-grpc_out=. telegram.proto
# выход в корень
cd ../../../..

# 5. Запуск сервера
go run cmd/server/main.go
```

---
### Конфигурация

| Переменная | Описание | Значение по умолчанию |
|------------|----------|---------------------|
| `GRPC_PORT` | Порт сервера | 50051 |
| `LOG_LEVEL` | Уровень логирования (debug/info/warn/error) | info |
| `MAX_SESSIONS` | Максимальное количество сессий | 10 |
| `TELEGRAM_APP_ID` | ID приложения Telegram (для реального режима) | - |
| `TELEGRAM_APP_HASH` | Hash приложения Telegram (для реального режима) | - |
---

### Реальное подключение к Telegram
Для работы с реальным Telegram необходимо:

1. Зарегистрировать приложение на my.telegram.org
2. Получить api_id и api_hash
3. Установить переменные окружения:

```bash
# можно указать в .env или вручную:
export TELEGRAM_APP_ID=ваш_app_id
export TELEGRAM_APP_HASH=ваш_app_hash
go run cmd/server/main.go
```
---
### Важно: Почему используется mock-режим?
В текущей версии сервиса используется mock-режим, так как:

Регистрация на my.telegram.org требует доступа к сайту, который может недоступен в РФ

1. API ключи необходимо получать индивидуально для каждого разработчика
2. Для демонстрации архитектуры и функциональности mock-режим полностью достаточен
3. При наличии валидных ключей сервис автоматически переключается в реальный режим работы с Telegram.
---
### gRPC Методы

| Метод | Описание | Request | Response |
|-------|----------|---------|----------|
| CreateSession | Создает новую сессию | {} | session_id, qr_code |
| SendMessage | Отправляет текстовое сообщение | session_id, peer, text | message_id |
| SubscribeMessages | Подписка на входящие сообщения (stream) | session_id | message_id, from, text, timestamp |
| DeleteSession | Удаляет сессию | session_id | {} |

---

### CreateSession

| Поле | Тип | Описание |
|------|-----|----------|
| session_id | string | Уникальный идентификатор сессии |
| qr_code | string | Данные для генерации QR кода |

Пример:
```bash
{
  "sessionId": "550e8400-e29b-41d4-a716-446655440000",
  "qrCode": "tg://login?token=mock_token_for_testing"
}
```
---

### SendMessage

| Поле | Тип | Описание |
|------|-----|----------|
| session_id | string | ID сессии |
| peer | string | Получатель (@username или @me) |
| text | string | Текст сообщения |
| message_id | int64 | ID отправленного сообщения |


Пример:
```bash
{
  "messageId": 1742954300123456789
}
```
---

### SubscribeMessages

| Поле | Тип | Описание |
|------|-----|----------|
| message_id | int64 | ID сообщения |
| from | string | Отправитель |
| text | string | Текст сообщения |
| timestamp | int64 | Время отправки (Unix timestamp) |


Пример (stream):
```bash
{
  "messageId": 1,
  "from": "test_bot",
  "text": "Test message #1 from session 550e8400-e29b-41d4-a716-446655440000",
  "timestamp": 1742954222
}
```
---

### DeleteSession

| Поле | Тип | Описание |
|------|-----|----------|
| - | - | Возвращает пустой объект |

```bash
Пример:
{}
```
---

## Тестирование

### grpcurl

| Метод | Команда |
|-------|---------|
| CreateSession | grpcurl -plaintext -d '{}' localhost:50051 telegram.TelegramService/CreateSession |
| SendMessage | grpcurl -plaintext -d '{"session_id":"xxx","peer":"@me","text":"Hello"}' localhost:50051 telegram.TelegramService/SendMessage |
| SubscribeMessages | grpcurl -plaintext -d '{"session_id":"xxx"}' localhost:50051 telegram.TelegramService/SubscribeMessages |
| DeleteSession | grpcurl -plaintext -d '{"session_id":"xxx"}' localhost:50051 telegram.TelegramService/DeleteSession |

---

### Postman

Шаг 1: Откройте Postman <br>
Шаг 2: Нажмите New -> gRPC Request <br>
Шаг 3: Введите URL: localhost:50051 <br>

| Метод | Request Body |
|-------|--------------|
| CreateSession | {} |
| SendMessage | {"session_id":"xxx","peer":"@me","text":"Hello"} |
| SubscribeMessages | {"session_id":"xxx"} |
| DeleteSession | {"session_id":"xxx"} |

---

### Ожидаемые результаты

| Тест | Ожидание |
|------|----------|
| CreateSession | sessionId + qrCode |
| SendMessage | messageId |
| SubscribeMessages | Stream сообщений каждые 5 секунд |
| DeleteSession | {} |
| Лимит сессий | Ошибка на 11-й сессии |
| Graceful shutdown | Сервер завершается корректно |

---

### Лимит сессий

| Действие | Результат |
|----------|-----------|
| Создание до 10 сессий | Успешно |
| Создание 11-й сессии | max sessions limit reached: 10 |
