package handlers

import (
	"task_scheduler/internal/entities"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) AddTask(newTask entities.Task) (string, error) {
	args := m.Called(newTask)
	return args.String(0), args.Error(1)
}
func (m *MockService) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) EditTask(updatedTask entities.Task) error {
	args := m.Called(updatedTask)
	return args.Error(0)
}

func (m *MockService) GetNextDate(now time.Time, date string, repeat string) (string, error) {
	args := m.Called(now, date, repeat)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetTask(id string) (entities.Task, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockService) GetTasks(target string) ([]entities.Task, error) {
	args := m.Called(target)
	return args.Get(0).([]entities.Task), args.Error(1)
}

type AuthService struct {
	mock.Mock
}

func (m *AuthService) GetJWT(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
