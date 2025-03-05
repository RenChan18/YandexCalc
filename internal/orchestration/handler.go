package orchestration

import (
	"encoding/json"
//	"errors"
//	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"calc_service/internal/calculator"
)

type Expression struct {
	ID         int     `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"` // pending, processing, completed, error
	Result     float64 `json:"result"`
	Tasks      []*Task `json:"-"`
}

type Task struct {
	ID            int     `json:"id"`
	ExpressionID  int     `json:"-"`
	Expression    string  `json:"expression,omitempty"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
	Result        float64 `json:"-"`
	Executed      bool    `json:"-"`
}

var (
	expressions = make(map[int]*Expression)
	tasksQueue  = make([]*Task, 0)
	exprMutex   sync.Mutex
	taskMutex   sync.Mutex
	nextExprID  = 1
	nextTaskID  = 1
)

// CalculateHandler обрабатывает POST /api/v1/calculate
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Expression string `json:"expression"`
	}
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || strings.TrimSpace(reqBody.Expression) == "" {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}

	// Проверяем выражение с помощью calculator.Calc
	if _, err := calculator.Calc(reqBody.Expression); err != nil {
		http.Error(w, "invalid expression: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Создаём новую запись выражения
	exprMutex.Lock()
	exprID := nextExprID
	nextExprID++
	expr := &Expression{
		ID:         exprID,
		Expression: reqBody.Expression,
		Status:     "pending",
	}
	expressions[exprID] = expr
	exprMutex.Unlock()

	// Определяем время операции по первому найденному оператору
	var opTime int
	var op string
	if strings.Contains(reqBody.Expression, "+") {
		op = "+"
		opTime, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	} else if strings.Contains(reqBody.Expression, "-") {
		op = "-"
		opTime, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	} else if strings.Contains(reqBody.Expression, "*") {
		op = "*"
		opTime, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	} else if strings.Contains(reqBody.Expression, "/") {
		op = "/"
		opTime, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	} else {
		op = "calc"
		opTime = 0
	}

	// Создаём задачу – для простоты задание содержит всё выражение
	taskMutex.Lock()
	taskID := nextTaskID
	nextTaskID++
	task := &Task{
		ID:            taskID,
		ExpressionID:  exprID,
		Expression:    reqBody.Expression,
		Operation:     op,
		OperationTime: opTime,
	}
	tasksQueue = append(tasksQueue, task)
	taskMutex.Unlock()

	// Привязываем задачу к выражению
	exprMutex.Lock()
	expr.Tasks = append(expr.Tasks, task)
	expr.Status = "processing"
	exprMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": exprID})
}

// GetExpressionsHandler обрабатывает GET /api/v1/expressions
func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	exprMutex.Lock()
	defer exprMutex.Unlock()

	exprs := make([]*Expression, 0, len(expressions))
	for _, e := range expressions {
		exprs = append(exprs, e)
	}

	resp := map[string]interface{}{
		"expressions": exprs,
	}
	json.NewEncoder(w).Encode(resp)
}

// GetExpressionByIDHandler обрабатывает GET /api/v1/expressions/{id}
func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	exprMutex.Lock()
	expr, ok := expressions[id]
	exprMutex.Unlock()
	if !ok {
		http.Error(w, "expression not found", http.StatusNotFound)
		return
	}

	resp := map[string]*Expression{"expression": expr}
	json.NewEncoder(w).Encode(resp)
}

// GetTaskHandler обрабатывает GET /internal/task
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	for _, t := range tasksQueue {
		if !t.Executed {
			resp := map[string]*Task{"task": t}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	http.Error(w, "no task", http.StatusNotFound)
}

// PostTaskResultHandler обрабатывает POST /internal/task
func PostTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}

	taskMutex.Lock()
	var task *Task
	for _, t := range tasksQueue {
		if t.ID == reqBody.ID {
			task = t
			break
		}
	}
	if task == nil {
		taskMutex.Unlock()
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	task.Executed = true
	task.Result = reqBody.Result
	taskMutex.Unlock()

	// Обновляем статус выражения
	exprMutex.Lock()
	if expr, ok := expressions[task.ExpressionID]; ok {
		allDone := true
		for _, t := range expr.Tasks {
			if !t.Executed {
				allDone = false
				break
			}
		}
		if allDone {
			expr.Result = reqBody.Result
			expr.Status = "completed"
		}
	}
	exprMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

