package services

import (
	"task_scheduler/internal/entities"
	"time"
)

type TaskServiceInterface interface {
	AddTask(newTask entities.Task) (string, error)
	DeleteTask(id string) error
	EditTask(updatedTask entities.Task) error
	GetNextDate(now time.Time, date string, repeat string) (string, error)
	GetTask(id string) (entities.Task, error)
	GetTasks(target string) ([]entities.Task, error)
}

type AuthServiceInterface interface {
	GetJWT(password string) (string, error)
}
