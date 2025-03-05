package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"calc_service/internal/calculator"
)

type Task struct {
	ID            int    `json:"id"`
	Expression    string `json:"expression,omitempty"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status: " + resp.Status)
	}
	var data struct {
		Task *Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Task, nil
}

func postTaskResult(taskID int, result float64) error {
	data := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	body, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to post result, status: " + resp.Status)
	}
	return nil
}

func worker(id int) {
	log.Printf("Worker %d запущен", id)
	for {
		task, err := getTask()
		if err != nil {
			log.Printf("Worker %d: ошибка получения задания: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if task == nil {
			// Нет заданий — ждём
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Worker %d: получено задание %d: %s", id, task.ID, task.Expression)
		// Имитируем задержку выполнения
		time.Sleep(time.Millisecond * time.Duration(task.OperationTime))
		// Вычисляем результат
		res, err := calculator.Calc(task.Expression)
		if err != nil {
			log.Printf("Worker %d: ошибка вычисления: %v", id, err)
			continue
		}
		// Отправляем результат обратно
		if err = postTaskResult(task.ID, res); err != nil {
			log.Printf("Worker %d: ошибка отправки результата: %v", id, err)
			continue
		}
		log.Printf("Worker %d: задание %d выполнено, результат = %v", id, task.ID, res)
	}
}

func main() {
	powerStr := os.Getenv("COMPUTING_POWER")
	if powerStr == "" {
		powerStr = "2"
	}
	power, err := strconv.Atoi(powerStr)
	if err != nil {
		power = 2
	}

	for i := 0; i < power; i++ {
		go worker(i + 1)
	}

	// Блокируем main, чтобы агент работал постоянно
	select {}
}

