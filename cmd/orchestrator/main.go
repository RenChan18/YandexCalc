package main

import (
    "log"
    "net/http"

    "calc_service/internal/orchestration"
)

func main() {
    mux := http.NewServeMux()

    // API для клиентов
    mux.HandleFunc("/api/v1/calculate", orchestration.AddExpressionHandler)
    // Дополнительно можно добавить GET /api/v1/expressions и GET /api/v1/expressions/{id}
    // Внутренние API для агентов
    mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            orchestration.GetTaskHandler(w, r)
        } else if r.Method == http.MethodPost {
            orchestration.SubmitTaskHandler(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    log.Println("Orchestrator server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
