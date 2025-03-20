package services

import "task_scheduler/internal/storage"

type TaskService struct {
	store storage.StorageInterface
}

func GetTaskService(store storage.StorageInterface) *TaskService {
	return &TaskService{store: store}
}
