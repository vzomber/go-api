package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	store, err := NewSQLiteTaskStore()
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
	mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler(store))
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler(store))
	mux.HandleFunc("POST /tasks", createTaskHandler(store))
	mux.HandleFunc("GET /tasks", getTasksHandler(store))
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("/", rootHandler)

	wrappedMux := loggingMiddleware(mux)
	err = http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
