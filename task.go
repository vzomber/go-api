package main

import (
	"encoding/json"
	"net/http"
)

var tasks = []Task{
	{ID: 1, Title: "Learn Go handlers", Done: false},
	{ID: 2, Title: "Build task API", Done: false},
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}
