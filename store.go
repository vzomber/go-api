package main

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type TaskFilter struct {
	Done   *bool
	Limit  *int
	Offset *int
}

type TaskStore interface {
	GetAll(filter TaskFilter) ([]Task, error)
	GetByID(id int) (Task, error)
	Add(title string) (Task, error)
	Update(id int, request UpdateTaskRequest) (Task, error)
	Delete(id int) (Task, error)
}

type SQLiteTaskStore struct {
	db *sql.DB
}

var ErrTaskNotFound = errors.New("task not found")

func OpenSQLiteTaskStore(path string) (*SQLiteTaskStore, error) {
	db, err := sql.Open("sqlite3", path)
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

	store, err := NewSQLiteTaskStore(db)
	if err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			log.Printf("failed to close db: %v", closeErr)
		}
		return nil, err
	}

	return store, nil
}

func NewSQLiteTaskStore(db *sql.DB) (*SQLiteTaskStore, error) {
	err := createTasksTable(db)
	if err != nil {
		return nil, err
	}

	return &SQLiteTaskStore{db: db}, nil
}

func (s *SQLiteTaskStore) seedData() error {
	var count int

	err := s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err = s.db.Exec(`
		INSERT INTO tasks (title, done)
		VALUES
			('Learn Go handlers', 0),
			('Build task API', 0)
	`)

	return err
}

func (s *SQLiteTaskStore) Close() error {
	return s.db.Close()
}

func createTasksTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0
		);
	`)
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

func appendIfMatchesFilter(tasks []Task, task Task, done *bool) []Task {
	if done != nil && *done != task.Done {
		return tasks
	}

	return append(tasks, task)
}

func applyPagination(tasks []Task, offset *int, limit *int) []Task {
	start := 0
	if offset != nil {
		start = *offset
	}

	if start > len(tasks) {
		return []Task{}
	}

	end := len(tasks)
	if limit != nil {
		end = start + *limit
	}

	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end]
}

func (s *SQLiteTaskStore) GetAll(taskFilter TaskFilter) ([]Task, error) {
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

		tasks = appendIfMatchesFilter(tasks, task, taskFilter.Done)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	tasks = applyPagination(tasks, taskFilter.Offset, taskFilter.Limit)

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

type MemoryTaskStore struct {
	tasks  []Task
	mu     sync.Mutex
	nextID int
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

func NewMemoryTaskStore(tasks []Task) *MemoryTaskStore {
	tasksCopy := make([]Task, len(tasks))
	copy(tasksCopy, tasks)

	return &MemoryTaskStore{
		tasks:  tasksCopy,
		nextID: calculateNextID(tasksCopy),
	}
}

func (s *MemoryTaskStore) Add(title string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		Title: title,
		ID:    s.nextID,
	}

	s.nextID++
	s.tasks = append(s.tasks, task)

	return task, nil
}

func (s *MemoryTaskStore) GetAll(taskFilter TaskFilter) ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []Task

	for _, task := range s.tasks {
		if taskFilter.Done != nil && task.Done != *taskFilter.Done {
			continue
		}

		tasks = append(tasks, task)
	}

	tasks = applyPagination(tasks, taskFilter.Offset, taskFilter.Limit)

	return tasks, nil
}

func (s *MemoryTaskStore) GetByID(requestedTaskID int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.ID == requestedTaskID {
			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}

func (s *MemoryTaskStore) Delete(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)

			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}

func (s *MemoryTaskStore) Update(id int, updateTaskBody UpdateTaskRequest) (Task, error) {
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

			return s.tasks[i], nil
		}
	}
	return Task{}, ErrTaskNotFound
}
