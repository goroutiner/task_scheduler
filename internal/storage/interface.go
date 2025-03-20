package storage

import (
	"task_scheduler/internal/entities"
)

type StorageInterface interface {
	PostTask(task entities.Task) (string, error)
	GetTasks() ([]entities.Task, error)
	SearchTasks(target string) ([]entities.Task, error)
	SearchTask(id string) (entities.Task, error)
	UpdateTask(task entities.Task) error
	DeleteTask(id string) error
}
