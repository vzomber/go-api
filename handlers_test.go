package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteTaskHandler(t *testing.T) {
	t.Run("deletes task", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
		})

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		tasks, err := store.GetAll(nil)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})
	t.Run("returns 404 when task not found", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{})

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/tasks/99", nil)

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rr.Code)
		}
	})
}

func TestUpdateTaskHandler(t *testing.T) {
	t.Run("updates task title", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "Old title", Done: false},
		})

		rr := httptest.NewRecorder()
		body := `{"title":"New title"}`
		req := httptest.NewRequest(http.MethodPatch, "/tasks/1", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var task Task

		err := json.NewDecoder(rr.Body).Decode(&task)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if task.Title != "New title" {
			t.Fatalf("expected title New title, got %q", task.Title)
		}

		if task.Done != false {
			t.Fatalf("expected done false, got %v", task.Done)
		}
	})

	t.Run("updates task", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "Old title", Done: false},
		})

		rr := httptest.NewRecorder()
		body := `{"title":"New title"}`
		req := httptest.NewRequest(http.MethodPatch, "/tasks/1", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var task Task

		err := json.NewDecoder(rr.Body).Decode(&task)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if task.ID != 1 {
			t.Fatalf("expected ID 1, got %d", task.ID)
		}

		if task.Title != "New title" {
			t.Fatalf("expected title New title, got %q", task.Title)
		}

		if task.Done != false {
			t.Fatalf("expected Done false, got %v", task.Done)
		}
	})

	t.Run("returns 404 when task not found", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "Old title", Done: false},
		})

		rr := httptest.NewRecorder()
		body := `{"title":"New title"}`
		req := httptest.NewRequest(http.MethodPatch, "/tasks/99", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rr.Code)
		}

		tasks, err := store.GetAll(nil)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		if tasks[0].Title != "Old title" {
			t.Fatalf("expected title unchanged, got %q", tasks[0].Title)
		}
	})
}

func TestCreateTaskHandler(t *testing.T) {
	t.Run("creates task", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{})

		rr := httptest.NewRecorder()
		body := `{"title":"Learn Js"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("POST /tasks", createTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rr.Code)
		}

		var task Task

		err := json.NewDecoder(rr.Body).Decode(&task)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if rr.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected application/json")
		}
		if task.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if task.Title != "Learn Js" {
			t.Fatalf("expected title Learn Js, got %q", task.Title)
		}

		tasks, err := store.GetAll(nil)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
	})
	t.Run("returns 400 for invalid json", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{})

		rr := httptest.NewRecorder()
		body := `{"title":New task}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("POST /tasks", createTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rr.Code)
		}

		if rr.Body.String() != "Invalid task data\n" {
			t.Fatalf("unexpected response body: %q", rr.Body.String())
		}

		tasks, err := store.GetAll(nil)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("returns 400 for invalid task data", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{})

		rr := httptest.NewRecorder()
		body := `{}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))

		mux := http.NewServeMux()
		mux.HandleFunc("POST /tasks", createTaskHandler(store))
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rr.Code)
		}

		if rr.Body.String() != "title cannot be empty\n" {
			t.Fatalf("unexpected response body: %q", rr.Body.String())
		}

		tasks, err := store.GetAll(nil)
		if err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})
}

func TestGetTaskHandler(t *testing.T) {
	t.Run("returns task", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
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
		store := NewMemoryTaskStore([]Task{
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
	store := NewMemoryTaskStore([]Task{
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
