package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
)

func main() {
    power := os.Getenv("COMPUTING_POWER")
    computingPower, err := strconv.Atoi(power)
    if err != nil || computingPower <= 0 {
        computingPower = 2 // значение по умолчанию
    }

    for i := 0; i < computingPower; i++ {
        go worker(i)
    }

    // Бесконечный цикл, чтобы агент не завершался
    select {}
}

func worker(id int) {
    log.Printf("Worker %d started", id)
    for {
        // Запрос задачи
        resp, err := http.Get("http://localhost:8080/internal/task")
        if err != nil {
            log.Printf("Worker %d: error fetching task: %v", id, err)
            time.Sleep(2 * time.Second)
            continue
        }
        if resp.StatusCode == http.StatusNotFound {
            // Нет задач, ждём некоторое время
            time.Sleep(1 * time.Second)
            continue
        }
        var taskResp struct {
            Task struct {
                ID            int    `json:"id"`
                Arg1          string `json:"arg1"`
                Arg2          string `json:"arg2"`
                Operation     string `json:"operation"`
                OperationTime int    `json:"operation_time"`
            } `json:"task"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
            log.Printf("Worker %d: error decoding task: %v", id, err)
            resp.Body.Close()
            continue
        }
        resp.Body.Close()

        // Имитация задержки выполнения
        time.Sleep(time.Duration(taskResp.Task.OperationTime) * time.Millisecond)

        // Простейшая логика вычисления (можно использовать internal/calculator)
        result := performOperation(taskResp.Task.Arg1, taskResp.Task.Arg2, taskResp.Task.Operation)

        // Отправка результата
        submitTaskResult(taskResp.Task.ID, result)
    }
}

func performOperation(arg1, arg2, operation string) float64 {
    // Здесь можно расширить логику или вызвать функцию из internal/calculator
    // Пример для сложения и умножения:
    var a, b float64
    a, _ = strconv.ParseFloat(arg1, 64)
    b, _ = strconv.ParseFloat(arg2, 64)
    switch operation {
    case "+":
        return a + b
    case "*":
        return a * b
    default:
        return 0
    }
}

func submitTaskResult(taskID int, result float64) {
    reqBody, _ := json.Marshal(map[string]interface{}{
        "id":     taskID,
        "result": result,
    })
    resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        log.Printf("Error submitting result for task %d: %v", taskID, err)
        return
    }
    resp.Body.Close()
}
