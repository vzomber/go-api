package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	if err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, World!")

	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func main() {
	mux := http.NewServeMux()
	store := newTaskStore([]Task{
		{ID: 1, Title: "Learn Go handlers", Done: false},
		{ID: 2, Title: "Build task API", Done: false},
	},
	)

	mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
	mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler(store))
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler(store))
	mux.HandleFunc("POST /tasks", createTaskHandler(store))
	mux.HandleFunc("GET /tasks", getTasksHandler(store))
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
