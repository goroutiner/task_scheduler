package services

import (
	"errors"
	"strconv"
	"task_scheduler/internal/entities"
	"time"
)

// AddTask добавляет задачу с параметрами, полученными из тела запроса.
func (s *TaskService) AddTask(newTask entities.Task) (string, error) {
	var (
		id  string
		err error
	)

	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	}

	if newTask.Title == "" {
		return "", errors.New("the task title is empty")
	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	date, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		return "", err
	}

	if date.Before(now) {
		if newTask.Repeat == "" {
			newTask.Date = now.Format("20060102")
		} else {
			newTask.Date, err = s.GetNextDate(now, newTask.Date, newTask.Repeat)
			if err != nil {
				return "", err
			}
		}
	}

	id, err = s.store.PostTask(newTask)

	return id, err
}

// editTask изменяет пармаетры задачи, полученные из тела запроса.
func (s *TaskService) EditTask(updatedTask entities.Task) error {
	var err error

	if updatedTask.Date == "" {
		updatedTask.Date = time.Now().Format("20060102")
	}

	if updatedTask.Title == "" {
		return errors.New("the task title is empty")
	}

	if _, err := strconv.Atoi(updatedTask.Id); err != nil {
		return errors.New("the id is not specified or is specified not correctly")
	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	date, err := time.Parse("20060102", updatedTask.Date)
	if err != nil {
		return err
	}

	if date.Before(now) {
		if updatedTask.Repeat == "" {
			updatedTask.Date = now.Format("20060102")
		} else {
			updatedTask.Date, err = s.GetNextDate(now, updatedTask.Date, updatedTask.Repeat)
			if err != nil {
				return err
			}
		}
	}

	err = s.store.UpdateTask(updatedTask)

	return err
}

// deleteTask удаляет задачу с id, полученным из параметра запроса.
func (s *TaskService) DeleteTask(id string) error {
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("the id is not specified or is specified not correctly")
	}

	err := s.store.DeleteTask(id)

	return err
}
