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

func insertTask(db *sql.DB, title string) (Task, error) {
	result, err := db.Exec(`INSERT INTO tasks(title) VALUES(?)`, title)
	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{ID: int(id), Title: title, Done: false}, nil
}

func updateTask(db *sql.DB, id int, title string, done bool) (Task, error) {
	result, err := db.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", title, done, id)
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

	return Task{ID: id, Title: title, Done: done}, nil
}

func deleteTask(db *sql.DB, id int) (Task, error) {
	task, err := getTaskById(db, id)
	if err != nil {
		return Task{}, err
	}

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
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
	log.Printf("all rows %v", rows)

	task, err := getTaskById(db, 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("task %d %v", 1, task)
	task, err = insertTask(db, "Sell a bike")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("inserted task %d %v", 1, task)

	task, err = updateTask(db, 3, "Build a company", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("updated task %d %v", 3, task)
}
