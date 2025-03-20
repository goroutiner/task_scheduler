package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"task_scheduler/internal/entities"
	"task_scheduler/internal/handlers"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestGetNextDate тестирует тестирует обработчик GetNextDate.
func TestGetNextDate(t *testing.T) {
	mockService := new(handlers.MockService)

	testNow := "20231010"
	testDate := "20231015"
	testRepeat := "d 5"
	expectedNextDate := "20231015"

	baseURL := "/api/nextdate"
	path := url.Values{}
	path.Add("now", testNow)
	path.Add("date", testDate)
	path.Add("repeat", testRepeat)
	fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())

	mux := http.NewServeMux()
	mux.HandleFunc("/api/nextdate", handlers.GetNextDate(mockService))

	t.Run("successful get next date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fullPath, nil)

		respRec := httptest.NewRecorder()

		mockService.On("GetNextDate", mock.Anything, mock.Anything, mock.Anything).Return("20231015", nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var actualNextDate string
		err := json.NewDecoder(respRec.Body).Decode(&actualNextDate)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		require.Equal(t, expectedNextDate, actualNextDate)
	})

	t.Run("valid error", func(t *testing.T) {
		testRepeat = "invalid format"
		path.Set("repeat", testRepeat)
		fullPath = fmt.Sprintf("%s?%s", baseURL, path.Encode())

		req := httptest.NewRequest(http.MethodGet, fullPath, nil)

		respRec := httptest.NewRecorder()

		mockService.ExpectedCalls = nil
		mockService.On("GetNextDate", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("some error"))
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusBadRequest, respRec.Code, "Ожидался статус 400, но получен %d", respRec.Code)

		var actualRes entities.Result
		expectedErrStr := "some error"

		err := json.NewDecoder(respRec.Body).Decode(&actualRes)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		require.Equal(t, expectedErrStr, actualRes.Error)
	})
}

// TestGetTasks тестирует обработчик GetTasks.
func TestGetTasks(t *testing.T) {
	mockService := new(handlers.MockService)

	tasks := []entities.Task{
		{
			Date:    "20230510",
			Title:   "Гонки на кольцевых",
			Comment: "надо отдохнуть",
			Repeat:  "",
		},
		{
			Date:    "20231011",
			Title:   "Сходить в кино",
			Comment: "взять попить",
			Repeat:  "",
		},
	}

	baseURL := "/api/tasks"
	path := url.Values{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", handlers.GetTasks(mockService))

	t.Run("successful search tasks", func(t *testing.T) {
		path.Add("search", "Гонки на кольцевых")
		fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())
		req := httptest.NewRequest(http.MethodGet, fullPath, nil)

		respRec := httptest.NewRecorder()

		mockService.On("GetTasks", mock.Anything).Return(tasks[:1], nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var response entities.Result

		err := json.NewDecoder(respRec.Body).Decode(&response)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		actualTasks := response.Tasks

		require.Equal(t, tasks[:1], actualTasks)
	})

	t.Run("successful get all tasks", func(t *testing.T) {
		path.Del("search")
		fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())
		req := httptest.NewRequest(http.MethodGet, fullPath, nil)

		respRec := httptest.NewRecorder()

		mockService.ExpectedCalls = nil
		mockService.On("GetTasks", mock.Anything).Return(tasks, nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var response entities.Result

		err := json.NewDecoder(respRec.Body).Decode(&response)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		actualTasks := response.Tasks

		require.Equal(t, tasks, actualTasks)
	})

	t.Run("valid error", func(t *testing.T) {
		fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())
		req := httptest.NewRequest(http.MethodGet, fullPath, nil)

		respRec := httptest.NewRecorder()

		mockService.ExpectedCalls = nil
		mockService.On("GetTasks", mock.Anything).Return([]entities.Task{}, errors.New("some error"))
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusInternalServerError, respRec.Code, "Ожидался статус 500, но получен %d", respRec.Code)

		var response entities.Result
		expectedErrStr := "some error"

		err := json.NewDecoder(respRec.Body).Decode(&response)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		require.Equal(t, expectedErrStr, response.Error)
	})
}

// TestDoneTask тестирует обработчик DoneTask.
func TestDoneTask(t *testing.T) {
	mockService := new(handlers.MockService)

	task := entities.Task{
		Id:      "1",
		Date:    "20231021",
		Title:   "Сходить в бильярд",
		Comment: "взять кий",
		Repeat:  "",
	}

	baseURL := "/api/task/done"
	path := url.Values{}
	path.Add("id", task.Id)
	fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())
	req := httptest.NewRequest(http.MethodPost, fullPath, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/task/done", handlers.DoneTask(mockService))

	t.Run("successful completion task", func(t *testing.T) {
		respRec := httptest.NewRecorder()

		mockService.On("GetTask", mock.Anything).Return(task, nil)
		mockService.On("DeleteTask", mock.Anything).Return(nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var actualTask entities.Task
		expectedTask := entities.Task{}

		err := json.NewDecoder(respRec.Body).Decode(&actualTask)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		require.Equal(t, expectedTask, actualTask)
	})

	t.Run("valid error", func(t *testing.T) {
		respRec := httptest.NewRecorder()

		mockService.ExpectedCalls = nil
		mockService.On("GetTask", mock.Anything).Return(task, nil)
		mockService.On("DeleteTask", mock.Anything).Return(errors.New("some error"))
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusInternalServerError, respRec.Code, "Ожидался статус 500, но получен %d", respRec.Code)

		var response entities.Result
		expectedErrStr := "some error"

		err := json.NewDecoder(respRec.Body).Decode(&response)
		require.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

		require.Equal(t, expectedErrStr, response.Error)
	})
}

// TestUpdateTasks тестирует обработчик UpdateTasks.
func TestUpdateTasks(t *testing.T) {
	mockService := new(handlers.MockService)

	task := entities.Task{
		Id:      "1",
		Date:    "20231021",
		Title:   "Сходить в боулинг",
		Comment: "взять ботинки",
		Repeat:  "",
	}

	baseURL := "/api/task"
	path := url.Values{}
	path.Add("id", "1")
	fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())
	body, _ := json.Marshal(task)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/task", handlers.UpdateTasks(mockService))

	t.Run("successful add task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, baseURL, bytes.NewReader(body))
		respRec := httptest.NewRecorder()

		mockService.On("AddTask", mock.Anything).Return("1", nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var actualRes entities.Result
		expectedId := "1"

		err := json.NewDecoder(respRec.Body).Decode(&actualRes)
		require.NoError(t, err)

		require.Equal(t, expectedId, actualRes.Id)
	})

	t.Run("successful get task", func(t *testing.T) {
		fullPath := fmt.Sprintf("%s?%s", baseURL, path.Encode())

		req := httptest.NewRequest(http.MethodGet, fullPath, nil)
		respRec := httptest.NewRecorder()

		mockService.On("GetTask", mock.Anything).Return(task, nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		var actualTask entities.Task

		err := json.NewDecoder(respRec.Body).Decode(&actualTask)
		require.NoError(t, err)

		require.Equal(t, task, actualTask)
	})

	t.Run("successful edit task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, baseURL, bytes.NewReader(body))
		respRec := httptest.NewRecorder()

		mockService.On("EditTask", mock.Anything).Return(nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		actualTask := entities.Task{}

		err := json.NewDecoder(respRec.Body).Decode(&actualTask)
		require.NoError(t, err)

		require.Equal(t, task, actualTask)
	})

	t.Run("successful delete task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fullPath, nil)
		respRec := httptest.NewRecorder()

		mockService.On("DeleteTask", mock.Anything).Return(nil)
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

		expectedTask := entities.Task{}
		actualTask := entities.Task{}

		err := json.NewDecoder(respRec.Body).Decode(&actualTask)
		require.NoError(t, err)

		require.Equal(t, expectedTask, actualTask)
	})

	t.Run("invalid request method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, baseURL, nil)
		respRec := httptest.NewRecorder()

		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusMethodNotAllowed, respRec.Code, "Ожидался статус 405, но получен %d", respRec.Code)
	})

	t.Run("valid error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fullPath, nil)
		respRec := httptest.NewRecorder()

		mockService.ExpectedCalls = nil
		mockService.On("DeleteTask", mock.Anything).Return(errors.New("some error"))
		mux.ServeHTTP(respRec, req)

		require.Equalf(t, http.StatusInternalServerError, respRec.Code, "Ожидался статус 500, но получен %d", respRec.Code)

		actualRes := entities.Result{}
		expectedErrStr := "some error"

		err := json.NewDecoder(respRec.Body).Decode(&actualRes)
		require.NoError(t, err)

		require.Equal(t, expectedErrStr, actualRes.Error)
	})
}

// TestAuthentication тестирует обработчик Authentication.
func TestAuthentication(t *testing.T) {
	mockAuth := new(handlers.AuthService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/signin", handlers.Authentication(mockAuth))
	t.Run("successful authentication", func(t *testing.T) {

		password := "qwerty12345678"
		body := map[string]string{"password": password}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewReader(bodyBytes))
		respRec := httptest.NewRecorder()

		mockAuth.On("GetJWT", mock.Anything).Return("some.test.token", nil)
		mux.ServeHTTP(respRec, req)

		require.Equal(t, http.StatusOK, respRec.Code)

		var result entities.Result
		expectedToken := "some.test.token"

		err := json.NewDecoder(respRec.Body).Decode(&result)

		require.NoError(t, err)
		require.Equal(t, expectedToken, result.Token)
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "wrong_password"
		body := map[string]string{"password": password}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewReader(bodyBytes))
		respRec := httptest.NewRecorder()

		mockAuth.ExpectedCalls = nil
		mockAuth.On("GetJWT", mock.Anything).Return("", errors.New("some error"))
		mux.ServeHTTP(respRec, req)

		require.Equal(t, http.StatusBadRequest, respRec.Code)

		var result entities.Result
		expectedError := "some error"

		err := json.NewDecoder(respRec.Body).Decode(&result)

		require.NoError(t, err)
		require.Empty(t, result.Token)
		require.Equal(t, expectedError, result.Error)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewReader([]byte("invalid body")))
		respRec := httptest.NewRecorder()

		mux.ServeHTTP(respRec, req)

		require.Equal(t, http.StatusBadRequest, respRec.Code)

		var result entities.Result

		err := json.NewDecoder(respRec.Body).Decode(&result)

		require.NoError(t, err)
		require.Empty(t, result.Token)
		require.NotEmpty(t, result.Error)
	})
}
