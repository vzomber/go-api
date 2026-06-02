package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func createTaskHandler(store *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newTask Task

		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, "Invalid task data", http.StatusBadRequest)
			return
		}

		createdTask := store.Add(newTask)

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdTask)
		if err != nil {
			http.Error(w, "Failed to encode task", http.StatusInternalServerError)
			return
		}
	}
}

func updateTaskHandler(store *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		requestedTaskID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var updateRequest UpdateTaskRequest

		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			log.Println("Invalid task data:", err)
			http.Error(w, "Invalid task data", http.StatusBadRequest)
			return
		}

		if updateRequest.Done == nil && updateRequest.Title == nil {
			log.Println("No fields to update:")
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		task, err := store.Update(requestedTaskID, updateRequest)
		if err != nil {
			log.Println("Failed to update task:", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			log.Println("Failed to encode task:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func deleteTaskHandler(store *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		requestedTaskID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, err := store.Delete(requestedTaskID)
		if err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			log.Println("Failed to delete task:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Failed to encode task:", err)
			return
		}
	}
}

func getTaskHandler(store *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		requestedTaskID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, err := store.GetByID(requestedTaskID)
		if err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			log.Println("Failed to get task by ID:", err)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("failed to encode task:", err)
			return
		}
	}
}

func getTasksHandler(store *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tasks := store.GetAll()

		err := json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
