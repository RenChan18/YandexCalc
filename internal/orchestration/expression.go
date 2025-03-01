package orchestration

import "time"

// ExpressionStatus описывает состояние вычисления выражения.
type ExpressionStatus string

const (
    StatusPending   ExpressionStatus = "pending"
    StatusRunning   ExpressionStatus = "running"
    StatusCompleted ExpressionStatus = "completed"
    StatusError     ExpressionStatus = "error"
)

// Expression представляет арифметическое выражение, принятое на обработку.
type Expression struct {
    ID        string            `json:"id"`
    Raw       string            `json:"raw"`
    Status    ExpressionStatus  `json:"status"`
    Result    string            `json:"result,omitempty"`
    CreatedAt time.Time         `json:"created_at"`
    Tasks     []Task            `json:"tasks,omitempty"`
}


// Task представляет отдельную операцию (например, сложение, умножение) из выражения.
type Task struct {
    ID            int     `json:"id"`
    Arg1          string  `json:"arg1"`
    Arg2          string  `json:"arg2"`
    Operation     string  `json:"operation"`
    OperationTime int     `json:"operation_time"` // В миллисекундах
    Result        float64 `json:"result,omitempty"`
    Done          bool    `json:"done"`
}
