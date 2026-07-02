# Task API

A simple REST API built in Go for learning backend development.

## Features

- CRUD operations for tasks
- SQLite persistence
- In-memory storage for tests
- Storage abstraction via `TaskStore` interface
- Filtering by completion status
- Pagination
- Search by title
- Unit tests for stores and handlers

## Tech Stack

- Go
- SQLite
- `net/http`
- `testing`
- `httptest`

## Running the project

Clone the repository:

```bash
git clone <repository-url>
cd <repository-name>
```

Run the API:

```bash
go run .
```

The server starts on:

```
http://localhost:8080
```

## Running tests

Run all tests:

```bash
go test ./...
```

## API Endpoints

### Get all tasks

```http
GET /tasks
```

Optional query parameters:

```http
GET /tasks?done=true
GET /tasks?done=false
GET /tasks?limit=10
GET /tasks?offset=5
GET /tasks?search=milk
```

Query parameters can be combined:

```http
GET /tasks?done=false&search=milk&limit=5&offset=0
```

### Get task by ID

```http
GET /tasks/{id}
```

### Create task

```http
POST /tasks
```

Example body:

```json
{
  "title": "Buy groceries"
}
```

### Update task

```http
PATCH /tasks/{id}
```

Example body:

```json
{
  "title": "Buy milk",
  "done": true
}
```

### Delete task

```http
DELETE /tasks/{id}
```

## Project Structure

```text
.
├── main.go
├── handlers.go
├── task_handlers.go
├── task.go
├── store.go
├── sqlite_store_test.go
├── memory_store_test.go
└── handlers_test.go
```

## Notes

This project was built as a learning exercise to practice Go, REST API design, storage abstraction, testing, and SQLite integration.
