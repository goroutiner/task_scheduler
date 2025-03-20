package services

import (
	"errors"
	"strconv"
	"task_scheduler/internal/entities"
)

// GetTasks возвращает задачи, содержащие строку или подсторку,
// полученную из параметров запроса.
func (s *TaskService) GetTasks(target string) ([]entities.Task, error) {
	var (
		tasks = []entities.Task{}
		err   error
	)

	// Поиск записи по значению параметра search, если он указан,
	// в противном случае показываются все задачи.
	if target != "" {
		tasks, err = s.store.SearchTasks(target)
	} else {
		tasks, err = s.store.GetTasks()
	}

	return tasks, err
}

// GetTask возвращает задачу по id, полученному из параметра запроса.
func (s *TaskService) GetTask(id string) (entities.Task, error) {
	var (
		task = entities.Task{}
		err  error
	)

	if _, err = strconv.Atoi(id); err != nil {
		return task, errors.New("the id is not specified or is specified not correctly")
	}

	task, err = s.store.SearchTask(id)

	return task, err
}
