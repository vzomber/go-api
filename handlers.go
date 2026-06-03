package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	if err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, World!")

	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}
