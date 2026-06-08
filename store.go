package main

import (
	"errors"
	"sync"
)

type TaskStore struct {
	tasks  []Task
	mu     sync.Mutex
	nextID int
}

var ErrTaskNotFound = errors.New("task not found")

func NewTaskStore(initialTasks []Task) *TaskStore {
	return &TaskStore{
		tasks:  initialTasks,
		nextID: len(initialTasks) + 1,
	}
}

func (s *TaskStore) Add(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		Title: title,
		ID:    s.nextID,
		Done:  false,
	}

	s.nextID++
	s.tasks = append(s.tasks, task)

	return task
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

			return s.tasks[i], nil
		}
	}
	return Task{}, ErrTaskNotFound
}
