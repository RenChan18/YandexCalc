# Calc Service

Этот проект представляет собой веб-сервис для вычисления арифметических выражений на Go.

## Структура проекта
```go
.
├── cmd
│   └── calc_service
│       └── main.go
├── go.mod
├── internal
│   └── calculator
│       └── calculator.go
├── README.md
└── tests
    └── calculator_test.go
```
- `cmd/calc_service/main.go`: Основной файл, запускающий веб-сервис.
- `internal/calculator/calculator.go`: Логика для обработки арифметических выражений.
- `tests/calculator_test.go`: Тесты для калькулятора.


## Запуск проекта

```bash
go run ./cmd/calc_service/...
```

## API

**POST /api/v1/calculate**

- **Тело запроса (JSON):**

| Параметр    | Тип   | Описание                            |
|-------------|-------|-------------------------------------|
| expression  | string| Арифметическое выражение для вычисления |

- **Ответ (JSON):**

| Параметр    | Тип   | Описание                            |
|-------------|-------|-------------------------------------|
| result      | string| Результат вычисления                |
| error       | string| Сообщение об ошибке                 |

**Пример запроса:**
```json
{
  "expression": "2+2*2"
}
```

## Возможные сценарии
#### Код 500
```bash
curl -i --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "error"
}'
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 20 Dec 2024 15:01:42 GMT
Content-Length: 35
```
ответ
```bash
{"error": "Internal server error"}

```
#### Код 200
```bash
curl -i --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
  "expression": "2*2+2"
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 20 Dec 2024 15:02:01 GMT
Content-Length: 22
```
ответ
```bash

{"result":"6.000000"}

```
#### Код 422
```bash
curl -i --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{
  "expression": "2*2+avc"
}'
HTTP/1.1 422 Unprocessable Entity
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 20 Dec 2024 15:02:10 GMT
Content-Length: 37
```
ответ
```bash

{"error": "Expression is not valid"}
```

