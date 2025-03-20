package services

import (
	"task_scheduler/internal/entities"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) PostTask(task entities.Task) (string, error) {
	args := m.Called(task)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetTasks() ([]entities.Task, error) {
	args := m.Called()
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockStorage) SearchTasks(target string) ([]entities.Task, error) {
	args := m.Called(target)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockStorage) SearchTask(id string) (entities.Task, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockStorage) UpdateTask(task entities.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockStorage) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
