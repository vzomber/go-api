package main

import (
	"errors"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (r CreateTaskRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title cannot be empty")
	}

	return nil
}

func (r UpdateTaskRequest) Validate() error {
	if r.Title != nil && strings.TrimSpace(*r.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if r.Done == nil && r.Title == nil {
		return errors.New("at least one field is required")
	}
	return nil
}
