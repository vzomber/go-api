package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "tasks.json")

	store := NewTaskStore([]Task{
		{ID: 1, Title: "Test task", Done: false},
	})
	store.filePath = filePath

	err := store.save()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected file to exist, got %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected file to contain data")
	}
}

func TestLoadExistingTasks(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "tasks.json")

	jsonData := `[
		{"id":1,"title":"First task","done":false},
		{"id":2,"title":"Second task","done":true}
	]`

	err := os.WriteFile(filePath, []byte(jsonData), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	store := NewTaskStore(nil)
	store.filePath = filePath

	err = store.load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(store.tasks))
	}
	if store.tasks[0].ID != 1 {
		t.Fatalf("expected first task ID 1, got %d", store.tasks[0].ID)
	}
	if store.tasks[1].Done != true {
		t.Fatal("expected second task to be done")
	}
}

func TestLoadMissingFileStartsEmpty(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "missing.json")

	store := NewTaskStore(nil)
	store.filePath = filePath

	err := store.load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(store.tasks))
	}
	if store.nextID != 1 {
		t.Fatalf("expected nextID to be 1, got %d", store.nextID)
	}
}

func TestNextIDCalculatedCorrectlyAfterLoad(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "tasks.json")

	jsonData := `[
		{"id":1,"title":"First","done":false},
		{"id":2,"title":"Second","done":false},
		{"id":5,"title":"Third","done":true}
	]`

	err := os.WriteFile(filePath, []byte(jsonData), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	store := NewTaskStore(nil)
	store.filePath = filePath

	err = store.load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store.nextID != 6 {
		t.Fatalf("expected nextID 6, got %d", store.nextID)
	}
}
