package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"calc_service/internal/calculator"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusUnprocessableEntity)
		return
	}

	// artificially generate error 500, if expression is "error"
	if req.Expression == "error" {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	result, err := calculator.Calc(req.Expression)
	if err != nil {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	resp := Response{Result: fmt.Sprintf("%f", result)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Server run on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}


