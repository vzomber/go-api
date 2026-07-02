package main

import (
	"database/sql"
	"errors"
	"testing"
)

func setupSQLiteTestStore(t *testing.T) *SQLiteTaskStore {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	store, err := NewSQLiteTaskStore(db)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close test db: %v", err)
		}
	})

	return store
}

func TestSQLiteTaskStoreAdd(t *testing.T) {
	store := setupSQLiteTestStore(t)

	task, err := store.Add("Learn SQLite")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID == 0 {
		t.Fatal("expected task ID to be set")
	}

	if task.Title != "Learn SQLite" {
		t.Fatalf("expected title %q, got %q", "Learn SQLite", task.Title)
	}

	if task.Done != false {
		t.Fatalf("expected Done false, got %v", task.Done)
	}
}

func TestSQLiteTaskStoreGetAll(t *testing.T) {
	store := setupSQLiteTestStore(t)

	_, err := store.Add("Get milk")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Add("Sell PS5")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tasks, err := store.GetAll(TaskFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	done := true
	updateBody := UpdateTaskRequest{Done: &done}

	_, err = store.Update(1, updateBody)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tasks, err = store.GetAll(TaskFilter{Done: &done})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestSQLiteTaskStoreGetByID(t *testing.T) {
	store := setupSQLiteTestStore(t)

	createdTask, err := store.Add("Test task")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	task, err := store.GetByID(createdTask.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Title != "Test task" {
		t.Fatalf("expected title %q, got %q", "Test task", task.Title)
	}
	if task.ID != createdTask.ID {
		t.Fatalf("expected ID %d, got %d", createdTask.ID, task.ID)
	}
}

func TestSQLiteTaskStoreUpdate(t *testing.T) {
	store := setupSQLiteTestStore(t)

	createdTask, err := store.Add("Old title")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	newTitle := "New title"
	done := true

	updatedTask, err := store.Update(createdTask.ID, UpdateTaskRequest{
		Title: &newTitle,
		Done:  &done,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedTask.ID != createdTask.ID {
		t.Fatalf("expected ID %d, got %d", createdTask.ID, updatedTask.ID)
	}

	if updatedTask.Title != newTitle {
		t.Fatalf("expected title %q, got %q", newTitle, updatedTask.Title)
	}

	if updatedTask.Done != done {
		t.Fatalf("expected done %v, got %v", done, updatedTask.Done)
	}
}

func TestSQLiteTaskStoreDelete(t *testing.T) {
	store := setupSQLiteTestStore(t)

	createdTask, err := store.Add("Task to delete")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	deletedTask, err := store.Delete(createdTask.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if deletedTask.ID != createdTask.ID {
		t.Fatalf("expected ID %d, got %d", createdTask.ID, deletedTask.ID)
	}

	if deletedTask.Title != createdTask.Title {
		t.Fatalf("expected title %q, got %q", createdTask.Title, deletedTask.Title)
	}

	_, err = store.GetByID(createdTask.ID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestSQLiteTaskStoreDeleteReturnsErrTaskNotFound(t *testing.T) {
	store := setupSQLiteTestStore(t)

	_, err := store.Delete(999)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestSQLiteTaskStoreGetAllWithPagination(t *testing.T) {
	store := setupSQLiteTestStore(t)

	_, err := store.Add("First task")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Add("Second task")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Add("Third task")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("applies limit", func(t *testing.T) {
		limit := 2

		tasks, err := store.GetAll(TaskFilter{Limit: &limit})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}

		if tasks[0].ID != 1 || tasks[1].ID != 2 {
			t.Fatalf("expected first two tasks, got %+v", tasks)
		}
	})

	t.Run("applies offset", func(t *testing.T) {
		offset := 1

		tasks, err := store.GetAll(TaskFilter{Offset: &offset})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}

		if tasks[0].ID != 2 || tasks[1].ID != 3 {
			t.Fatalf("expected second and third tasks, got %+v", tasks)
		}
	})

	t.Run("applies offset and limit", func(t *testing.T) {
		offset := 1
		limit := 1

		tasks, err := store.GetAll(TaskFilter{
			Offset: &offset,
			Limit:  &limit,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		if tasks[0].ID != 2 {
			t.Fatalf("expected task ID 2, got %d", tasks[0].ID)
		}
	})
}

func TestSQLiteTaskStoreGetAllSearch(t *testing.T) {
	store := setupSQLiteTestStore(t)

	_, err := store.Add("Get milk")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Add("Sell PS5")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Add("Buy MILK chocolate")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("matches by title", func(t *testing.T) {
		search := "milk"

		tasks, err := store.GetAll(TaskFilter{Search: &search})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("is case-insensitive", func(t *testing.T) {
		search := "MILK"

		tasks, err := store.GetAll(TaskFilter{Search: &search})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("returns empty when no match", func(t *testing.T) {
		search := "unknown"

		tasks, err := store.GetAll(TaskFilter{Search: &search})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 0 {
			t.Fatalf("expected 0 tasks, got %d", len(tasks))
		}
	})
}
