package orchestration

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/google/uuid"
    "calc_service/internal/calculator" // импорт старой логики вычислений
)

var (
    // Здесь можно хранить выражения в памяти для простоты примера.
    expressions = make(map[string]*Expression)
    // Очередь задач для агентов.
    taskQueue = NewTaskQueue(100)
)

// AddExpressionHandler обрабатывает POST /api/v1/calculate.
func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Expression string `json:"expression"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
        http.Error(w, `{"error": "Invalid input"}`, http.StatusUnprocessableEntity)
        return
    }
	// Используем функцию из пакета calculator для проверки выражения.
	_, err := calculator.Calc(req.Expression)
	if err != nil {
	 http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
	 return
	}
	id := uuid.New().String()
    exp := &Expression{
        ID:        id,
        Raw:       req.Expression,
        Status:    StatusPending,
        CreatedAt: time.Now(),
    }

    // Пример разбиения выражения на задачи (псевдокод)
    // Можно вызвать функцию из internal/calculator для парсинга
    // Здесь для простоты предположим, что выражение "2*2+2" делится на 2 задачи:
    // 1. Вычислить 2*2
    // 2. Вычислить результат предыдущей операции + 2
    if req.Expression == "2*2+2" {
        task1 := Task{
            ID:            1,
            Arg1:          "2",
            Arg2:          "2",
            Operation:     "*",
            OperationTime: getOperationTime("multiplication"),
        }
        task2 := Task{
            ID:            2,
            Arg1:          "4", // ожидаем результат task1
            Arg2:          "2",
            Operation:     "+",
            OperationTime: getOperationTime("addition"),
        }
        exp.Tasks = []Task{task1, task2}

        // Добавляем задачи в очередь
        taskQueue.Enqueue(task1)
        taskQueue.Enqueue(task2)
    } else {
        // Можно попытаться вычислить синхронно или вернуть ошибку.
        http.Error(w, `{"error": "Expression format not supported"}`, http.StatusUnprocessableEntity)
        return
    }

    expressions[id] = exp

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getOperationTime(operation string) int {
    // Читаем переменные окружения, задающие время выполнения.
    var envKey string
    switch operation {
    case "addition":
        envKey = "TIME_ADDITION_MS"
    case "subtraction":
        envKey = "TIME_SUBTRACTION_MS"
    case "multiplication":
        envKey = "TIME_MULTIPLICATIONS_MS"
    case "division":
        envKey = "TIME_DIVISIONS_MS"
    default:
        return 0
    }
    msStr := os.Getenv(envKey)
    ms, err := strconv.Atoi(msStr)
    if err != nil {
        // Если переменная не установлена, возвращаем дефолтное значение, например, 50 мс.
        return 50
    }
    return ms
}

// GetTaskHandler обрабатывает GET /internal/task для агентов.
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
    task, ok := taskQueue.Dequeue()
    if !ok {
        http.Error(w, `{"error": "No tasks available"}`, http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]Task{"task": task})
}

// SubmitTaskHandler обрабатывает POST /internal/task для принятия результата.
func SubmitTaskHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID     int     `json:"id"`
        Result float64 `json:"result"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid input"}`, http.StatusUnprocessableEntity)
        return
    }
    // Здесь необходимо найти задачу и обновить её статус, а затем – обновить состояние выражения.
    // Для простоты примера просто логируем результат.
    log.Printf("Task %d completed with result %f", req.ID, req.Result)
    w.WriteHeader(http.StatusOK)
}
