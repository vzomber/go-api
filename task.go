package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type UpdateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newTask)
	if err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
		return
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateRequest UpdateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	requestedTaskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == requestedTaskID {

			if updateRequest.Done != nil {
				tasks[i].Done = *updateRequest.Done
			}

			if updateRequest.Title != nil {
				tasks[i].Title = *updateRequest.Title
			}

			err := json.NewEncoder(w).Encode(tasks[i])
			if err != nil {
				http.Error(w, "Cannot encode task", http.StatusInternalServerError)
				fmt.Println("failed to encode task:", err)
			}
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	requestedTaskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == requestedTaskID {

			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)

			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	requestedTaskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == requestedTaskID {
			err := json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, "Cannot encode task", http.StatusInternalServerError)
				fmt.Println("failed to encode task:", err)
			}
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}
