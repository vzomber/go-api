package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID    int
	Title string
	Done  bool
}

func getAllRows(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, done FROM tasks")
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

func createTasksTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		return err
	}

	log.Println("table created")
	return nil
}

func getTaskById(db *sql.DB, id int) (Task, error) {
	row := db.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id)

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

func main() {
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = createTasksTable(db)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := getAllRows(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(rows)

	task, err := getTaskById(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(task)
}
