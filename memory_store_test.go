package main

import (
	"errors"
	"testing"
)

func TestGetAll(t *testing.T) {
	t.Run("return all tasks", func(t *testing.T) {
		mockTasks := []Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		}
		store := NewMemoryTaskStore(mockTasks)

		tasks, err := store.GetAll(TaskFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}

		if tasks[0].Title != mockTasks[0].Title {
			t.Fatalf("expected first task title %q, got %q", mockTasks[0].Title, tasks[0].Title)
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("adds task and assigns ID", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{})

		createdTask, err := store.Add("To go shopping")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if createdTask.ID != 1 {
			t.Fatalf("expected ID to be %d, but got %d", 1, createdTask.ID)
		}

		if createdTask.Title != "To go shopping" {
			t.Fatalf("expected title %q, got %q", "To go shopping", createdTask.Title)
		}

		if createdTask.Done != false {
			t.Fatalf("expected Done to be false, got %v", createdTask.Done)
		}

		tasks, err := store.GetAll(TaskFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected len to be %d, but got %d", 1, len(tasks))
		}
	})
}

func TestGetByID(t *testing.T) {
	t.Run("returns task when task exists", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		task, err := store.GetByID(1)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if task.ID != 1 {
			t.Fatalf("expected ID 1, got %d", task.ID)
		}
		if task.Title != "First task" {
			t.Fatalf("expected title %q, got %q", "First task", task.Title)
		}
	})

	t.Run("returns ErrTaskNotFound when task does not exist", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		_, err := store.GetByID(999)
		if !errors.Is(err, ErrTaskNotFound) {
			t.Fatalf("expected ErrTaskNotFound, got %v", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("updates Title", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		title := "Third task"

		task, err := store.Update(1, UpdateTaskRequest{Title: &title})
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if task.Title != "Third task" {
			t.Fatalf("expected title %q, got %q", "Third task", task.Title)
		}
		if task.Done != false {
			t.Fatalf("expected Done to remain false")
		}
	})

	t.Run("updates Done", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		done := true

		task, err := store.Update(1, UpdateTaskRequest{
			Done: &done,
		})
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if task.Done != true {
			t.Fatalf("expected Done to be true")
		}
		if task.Title != "First task" {
			t.Fatalf("expected title to remain unchanged")
		}
	})

	t.Run("returns ErrTaskNotFound when task does not exist", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		_, err := store.Update(99, UpdateTaskRequest{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Fatalf("expected ErrTaskNotFound, got %v", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("returns deleted task", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		task, err := store.Delete(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.ID != 1 {
			t.Fatalf("expected ID 1, got %d", task.ID)
		}

		tasks, err := store.GetAll(TaskFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task remaining, got %d", len(tasks))
		}
		if tasks[0].ID != 2 {
			t.Fatalf("expected remaining task ID 2, got %d", tasks[0].ID)
		}
	})

	t.Run("returns ErrTaskNotFound when task does not exist", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		_, err := store.Delete(99)
		if !errors.Is(err, ErrTaskNotFound) {
			t.Fatalf("expected ErrTaskNotFound, got %v", err)
		}
	})
}

func TestTaskStoreFlow(t *testing.T) {
	t.Run("can add, delete and get tasts", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		tasks, err := store.GetAll(TaskFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		createdTask, err := store.Add("Third task")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if createdTask.ID != 3 {
			t.Fatalf("expected created task ID 3, got %d", createdTask.ID)
		}

		task, err := store.GetByID(3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "Third task" {
			t.Fatalf("expected title %q, got %q", "Third task", task.Title)
		}

		_, err = store.Delete(3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		tasks, err = store.GetAll(TaskFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks after delete, got %d", len(tasks))
		}
	})
}

func TestMemoryTaskStoreGetAllFiltersByDone(t *testing.T) {
	t.Run("filters with done = true", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		done := true
		tasks, err := store.GetAll(TaskFilter{Done: &done})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if !tasks[0].Done {
			t.Fatalf("expected task to be done")
		}
	})
	t.Run("filters with done = false", func(t *testing.T) {
		store := NewMemoryTaskStore([]Task{
			{ID: 1, Title: "First task", Done: false},
			{ID: 2, Title: "Second task", Done: true},
		})

		done := false
		tasks, err := store.GetAll(TaskFilter{Done: &done})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Done {
			t.Fatalf("expected task to be not done")
		}
	})
}
