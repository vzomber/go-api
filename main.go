package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	mux := http.NewServeMux()
	store, err := OpenSQLiteTaskStore("tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	err = store.seedData()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := store.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	mux.HandleFunc("PATCH /tasks/{id}", updateTaskHandler(store))
	mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler(store))
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler(store))
	mux.HandleFunc("POST /tasks", createTaskHandler(store))
	mux.HandleFunc("GET /tasks", getTasksHandler(store))
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("/", rootHandler)

	wrappedMux := loggingMiddleware(mux)
	err = http.ListenAndServe(":"+port, wrappedMux)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
