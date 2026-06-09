package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTaskHandler(t *testing.T) {
	t.Run("returns task", func(t *testing.T) {
		store := NewTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
		})

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		rr := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /tasks/{id}", getTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		var task Task

		err := json.NewDecoder(rr.Body).Decode(&task)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if task.ID != 1 {
			t.Fatalf("expected ID 1, got %d", task.ID)
		}
		if task.Title != "First task" {
			t.Fatalf("expected title First task, got %q", task.Title)
		}

	})
	t.Run("returns 404 when task not found", func(t *testing.T) {
		store := NewTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
		})

		req := httptest.NewRequest(http.MethodGet, "/tasks/99", nil)
		rr := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /tasks/{id}", getTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rr.Code)
		}
	})
}

func TestGetTasksHandler(t *testing.T) {
	store := NewTaskStore([]Task{
		{ID: 1, Title: "First task", Done: false},
		{ID: 2, Title: "Second task", Done: true},
	})

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler := getTasksHandler(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var tasks []Task

	err := json.NewDecoder(rr.Body).Decode(&tasks)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != 1 {
		t.Fatalf("expected first task ID 1, got %d", tasks[0].ID)
	}

}
