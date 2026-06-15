package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (r CreateTaskRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title cannot be empty")
	}

	return nil
}

func (r UpdateTaskRequest) Validate() error {
	if r.Title != nil && strings.TrimSpace(*r.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if r.Done == nil && r.Title == nil {
		return errors.New("at least one field is required")
	}
	return nil
}

func createTaskHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid task data", http.StatusBadRequest)
			return
		}

		err = request.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdTask, err := store.Add(request.Title)
		if err != nil {
			log.Printf("failed to save task: %v", err)
			http.Error(w, "Failed to save new task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdTask)
		if err != nil {
			http.Error(w, "Failed to encode task", http.StatusInternalServerError)
			return
		}
	}
}

func updateTaskHandler(store TaskStore) http.HandlerFunc {
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
			log.Println("Invalid task data: ", err)
			http.Error(w, "Invalid task data", http.StatusBadRequest)
			return
		}

		err = updateRequest.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updateRequest.Done == nil && updateRequest.Title == nil {
			log.Println("No fields to update")
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		task, err := store.Update(requestedTaskID, updateRequest)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				http.Error(w, "Task not found", http.StatusNotFound)
				log.Printf("task %d not found", requestedTaskID)
				return
			}

			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println("Failed to delete task: ", err)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			log.Println("Failed to encode task: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func deleteTaskHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		requestedTaskID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, err := store.Delete(requestedTaskID)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				http.Error(w, "Task not found", http.StatusNotFound)
				log.Printf("task %d not found", requestedTaskID)
				return
			}

			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println("Failed to delete task: ", err)
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

func getTaskHandler(store TaskStore) http.HandlerFunc {
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
			if errors.Is(err, ErrTaskNotFound) {
				http.Error(w, "Task not found", http.StatusNotFound)
				log.Printf("task %d not found", requestedTaskID)
				return
			}

			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Println("Failed to delete task: ", err)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("failed to encode task: ", err)
			return
		}
	}
}

func getTasksHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var doneFilter *bool

		doneParam := r.URL.Query().Get("done")
		if doneParam != "" {
			done, err := strconv.ParseBool(doneParam)
			if err != nil {
				log.Println("invalid done filter:", err)
				http.Error(w, "invalid done filter", http.StatusBadRequest)
				return
			}
			doneFilter = &done
		}

		tasks, err := store.GetAll(doneFilter)
		if err != nil {
			log.Println("could not get tasks:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
