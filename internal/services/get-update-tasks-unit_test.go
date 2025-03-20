package services_test

import (
	"errors"
	"fmt"
	"task_scheduler/internal/entities"
	"task_scheduler/internal/services"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	validTasksTableForUpdate = []entities.Task{
		{
			Id:      "1",
			Date:    "20240220",
			Title:   "Просмотр фильма",
			Comment: "с попкорном",
			Repeat:  "",
		},
		{
			Id:      "2",
			Date:    "20260701",
			Title:   "Сходить в бассейн",
			Comment: "",
			Repeat:  "",
		},
		{
			Id:      "3",
			Date:    "20250101",
			Title:   "Оплатить коммуналку",
			Comment: "",
			Repeat:  "d 30",
		},
		{
			Id:      "4",
			Date:    "20241231",
			Title:   "Поплавать",
			Comment: "Бассейн с тренером",
			Repeat:  "d 7",
		},
		{
			Id:      "5",
			Date:    "16890220",
			Title:   "Поплавать",
			Comment: "Бассейн с тренером",
			Repeat:  "d 7",
		},
		{
			Id:      "6",
			Date:    "",
			Title:   "Просмотр фильма",
			Comment: "с попкорном",
			Repeat:  "",
		},
	}

	invalidTasksTableForUpdate = []entities.Task{
		{
			Id:      "1",
			Date:    "20240220",
			Title:   "",
			Comment: "с попкорном",
			Repeat:  "",
		},
		{
			Id:      "2",
			Date:    "isnotdate",
			Title:   "Оплатить коммуналку",
			Comment: "",
			Repeat:  "d 30",
		},
		{
			Id:      "3",
			Date:    "20240507",
			Title:   "Встретится с Васей",
			Comment: "в 18:00",
			Repeat:  "invalidrepeat",
		},
		{
			Id:      "",
			Date:    "20240507",
			Title:   "Просмотр фильма",
			Comment: "с попкорном",
			Repeat:  "invalidrepeat",
		},
	}

	validTasksTableForGet = []entities.Task{
		{
			Id:      "1",
			Date:    "20240220",
			Title:   "Просмотр фильма",
			Comment: "с попкорном",
			Repeat:  "",
		},
		{
			Id:      "2",
			Date:    "20240229",
			Title:   "Просмотр фильма",
			Comment: "с сухариками",
			Repeat:  "",
		},
		{
			Id:      "2",
			Date:    "20250712",
			Title:   "Сходить в душ",
			Comment: "с мочалкой",
			Repeat:  "",
		},
	}
)

// TestAddTask тестирует метод AddTask сервиса задач.
func TestAddTask(t *testing.T) {
	mockStore := new(services.MockStorage)
	s := services.GetTaskService(mockStore)

	t.Run("post valid task", func(t *testing.T) {
		for i, newTask := range validTasksTableForUpdate {
			mockPostTask := mockStore.On("PostTask", mock.Anything).Return(fmt.Sprint(i+1), nil)
			id, err := s.AddTask(newTask)

			require.Equal(t, fmt.Sprint(i+1), id)
			require.NoError(t, err)

			mockPostTask.Unset()
			mockStore.AssertExpectations(t)
		}
	})

	t.Run("post invalid task", func(t *testing.T) {
		for _, newTask := range invalidTasksTableForUpdate {
			mockStore.On("PostTask", mock.Anything).Return("", nil)
			_, err := s.AddTask(newTask)

			require.Error(t, err)

			mockStore.AssertNotCalled(t, "PostTask")
		}
	})
}

// TestEditTask тестирует метод EditTask сервиса задач.
func TestEditTask(t *testing.T) {
	mockStore := new(services.MockStorage)
	s := services.GetTaskService(mockStore)

	t.Run("update valid task", func(t *testing.T) {
		for _, updatedTask := range validTasksTableForUpdate {
			mockStore.On("UpdateTask", mock.Anything).Return(nil)
			err := s.EditTask(updatedTask)

			require.NoError(t, err)

			mockStore.ExpectedCalls = nil
			mockStore.AssertExpectations(t)
		}
	})

	t.Run("update invalid task", func(t *testing.T) {
		for _, updatedTask := range invalidTasksTableForUpdate {
			mockStore.On("UpdateTask", mock.Anything).Return(nil)
			err := s.EditTask(updatedTask)

			require.Error(t, err)
			mockStore.AssertNotCalled(t, "PostTask")
		}
	})
}

// TestDeleteTask тестирует метод DeleteTask сервиса задач.
func TestDeleteTask(t *testing.T) {
	mockStore := new(services.MockStorage)
	s := services.GetTaskService(mockStore)

	mockStore.On("DeleteTask", mock.Anything).Return(nil)
	t.Run("delete valid task", func(t *testing.T) {
		testId := "1"

		err := s.DeleteTask(testId)

		require.NoError(t, err)
	})

	t.Run("delete valid task", func(t *testing.T) {
		testId := "isnotnum"

		err := s.DeleteTask(testId)

		require.Error(t, err)
		mockStore.AssertNotCalled(t, "DeleteTask")
	})
}

// TestGetTasks тестирует метод GetTasks сервиса задач.
func TestGetTasks(t *testing.T) {
	mockStore := new(services.MockStorage)
	s := services.GetTaskService(mockStore)

	t.Run("search valid tasks", func(t *testing.T) {
		testTarget := "Просмотр фильма"

		mockStore.On("SearchTasks", testTarget).Return(validTasksTableForGet[:3], nil)
		actualTasks, err := s.GetTasks(testTarget)

		require.NoError(t, err)
		require.Equal(t, validTasksTableForGet[:3], actualTasks)

	})

	t.Run("get valid tasks", func(t *testing.T) {
		testTarget := ""

		mockStore.On("GetTasks").Return(validTasksTableForGet, nil)
		actualTasks, err := s.GetTasks(testTarget)

		require.NoError(t, err)
		require.Equal(t, validTasksTableForGet, actualTasks)
	})

	t.Run("search invalid tasks", func(t *testing.T) {
		testTarget := "Просмотр матча"

		mockStore.ExpectedCalls = nil
		mockStore.On("SearchTasks", testTarget).Return([]entities.Task{}, errors.New("failed"))

		tasks, err := s.GetTasks(testTarget)

		require.Error(t, err)
		require.Empty(t, tasks)
	})

	t.Run("get invalid tasks", func(t *testing.T) {
		testTarget := ""

		mockStore.ExpectedCalls = nil
		mockStore.On("GetTasks").Return([]entities.Task{}, errors.New("failed"))

		tasks, err := s.GetTasks(testTarget)

		require.Error(t, err)
		require.Empty(t, tasks)
	})
}

// TestGetTask тестирует метод GetTask сервиса задач.
func TestGetTask(t *testing.T) {
	mockStore := new(services.MockStorage)
	s := services.GetTaskService(mockStore)

	t.Run("search valid task", func(t *testing.T) {
		testId := "1"

		mockStore.On("SearchTask", testId).Return(validTasksTableForGet[0], nil)
		actualTask, err := s.GetTask(testId)

		require.NoError(t, err)
		require.Equal(t, validTasksTableForGet[0], actualTask)
	})

	t.Run("search invalid task", func(t *testing.T) {
		testId := ""

		task, err := s.GetTask(testId)

		require.Error(t, err)
		require.Empty(t, task)
	})

}
