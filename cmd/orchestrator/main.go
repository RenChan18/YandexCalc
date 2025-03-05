package main

import (
	"log"
	"net/http"

	"calc_service/internal/orchestration"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestration.CalculateHandler)
	http.HandleFunc("/api/v1/expressions/", orchestration.GetExpressionByIDHandler) // для конкретного выражения
	http.HandleFunc("/api/v1/expressions", orchestration.GetExpressionsHandler)
	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orchestration.GetTaskHandler(w, r)
		case http.MethodPost:
			orchestration.PostTaskResultHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := "8080"
	log.Printf("Оркестратор запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

