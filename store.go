package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type TaskStore struct {
	tasks    []Task
	mu       sync.Mutex
	nextID   int
	filePath string
}

type SQLiteTaskStore struct {
	db *sql.DB
}

var ErrTaskNotFound = errors.New("task not found")

func NewSQLiteTaskStore() (*SQLiteTaskStore, error) {
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	store := &SQLiteTaskStore{
		db: db,
	}

	err = store.createTasksTable()
	if err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteTaskStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteTaskStore) createTasksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0
		)
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteTaskStore) Add(title string) (Task, error) {
	result, err := s.db.Exec(`INSERT INTO tasks(title) VALUES(?)`, title)
	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{ID: int(id), Title: title, Done: false}, nil
}

func (s *SQLiteTaskStore) GetAll() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close db cursor: %v", err)
		}
	}()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *SQLiteTaskStore) GetByID(id int) (Task, error) {
	row := s.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id)

	var task Task

	err := row.Scan(&task.ID, &task.Title, &task.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}

	return task, nil
}

func (s *SQLiteTaskStore) Delete(id int) (Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return Task{}, err
	}

	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Task{}, err
	}

	if rowsAffected == 0 {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *SQLiteTaskStore) Update(id int, request UpdateTaskRequest) (Task, error) {
	result, err := s.db.Exec(`
		UPDATE tasks SET
			title = COALESCE(?, title),
			done  = COALESCE(?, done)
		WHERE id = ?`,
		request.Title, request.Done, id)
	if err != nil {
		return Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Task{}, err
	}

	if rowsAffected == 0 {
		return Task{}, ErrTaskNotFound
	}

	return s.GetByID(id)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func NewTaskStoreFromFile(filePath string) (*TaskStore, error) {
	store := &TaskStore{
		filePath: filePath,
	}

	err := store.load()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *TaskStore) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.tasks = []Task{}
			s.nextID = 1
			return nil
		}

		return err
	}

	err = json.Unmarshal(data, &s.tasks)
	if err != nil {
		return err
	}

	s.nextID = calculateNextID(s.tasks)

	return nil
}

func calculateNextID(tasks []Task) int {
	maxID := 0

	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1
}

func (s *TaskStore) save() error {
	if s.filePath == "" {
		return nil
	}

	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewTaskStore(initialTasks []Task) *TaskStore {
	return &TaskStore{
		tasks:  initialTasks,
		nextID: calculateNextID(initialTasks),
	}
}

func (s *TaskStore) Add(title string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		Title: title,
		ID:    s.nextID,
		Done:  false,
	}

	s.nextID++
	s.tasks = append(s.tasks, task)

	err := s.save()
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *TaskStore) GetAll() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasksCopy := make([]Task, len(s.tasks))
	copy(tasksCopy, s.tasks)

	return tasksCopy
}

func (s *TaskStore) GetByID(requestedTaskID int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.ID == requestedTaskID {
			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}

func (s *TaskStore) Delete(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)

			err := s.save()
			if err != nil {
				return Task{}, err
			}

			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}

func (s *TaskStore) Update(id int, updateTaskBody UpdateTaskRequest) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {

			if updateTaskBody.Done != nil {
				s.tasks[i].Done = *updateTaskBody.Done
			}
			if updateTaskBody.Title != nil {
				s.tasks[i].Title = *updateTaskBody.Title
			}

			err := s.save()
			if err != nil {
				return Task{}, err
			}

			return s.tasks[i], nil
		}
	}
	return Task{}, ErrTaskNotFound
}
