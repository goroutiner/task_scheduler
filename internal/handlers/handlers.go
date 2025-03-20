package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"task_scheduler/internal/entities"
	"task_scheduler/internal/services"
	"time"
)

// Authentication получает пароль через HTTPS, и если он является валидным,
// то отправляет HTTP ответ, содержащий сгенерированный токен.
func Authentication(authService services.AuthServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			signedToken   string
			takenPassword map[string]string
			err           error
		)

		err = json.NewDecoder(r.Body).Decode(&takenPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		signedToken, err = authService.GetJWT(takenPassword["password"])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(entities.Result{Token: signedToken})
	}
}

// GetNextDate получает значения параметров now, date, repeat из параметров запроса и
// с их помощью возвращает HTTP ответ, содержащий следующую ближайшую дату.
func GetNextDate(s services.TaskServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")

		timeNow, err := time.Parse("20060102", now)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		res, err := s.GetNextDate(timeNow, date, repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(res)
	}
}

// GetTasks возвращает HTTP ответ, содержащий список всех существующих задач.
func GetTasks(s services.TaskServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err   error
			tasks []entities.Task
		)

		target := r.FormValue("search")
		tasks, err = s.GetTasks(target)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(entities.Result{Tasks: tasks})
	}
}

// DoneTask завершает или обновляет дату задачи, если поле repeat не пустое и
// возвращает пустой JSON в случае успешной обработки.
func DoneTask(s services.TaskServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			task entities.Task
			id   string
			err  error
		)

		if r.FormValue("id") == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.Result{Error: "id не указан или указан некорректно"})
			return
		}

		id = r.FormValue("id")
		task, err = s.GetTask(id)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		if task.Repeat != "" {
			nextDate, err := s.GetNextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
				return
			}

			task.Date = nextDate

			err = s.EditTask(task)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
				return
			}
		} else {
			err = s.DeleteTask(id)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(entities.Task{})
	}
}

// UpdateTasks обрабатывает несколько методов: POST, GET, PUT, DELETE.
func UpdateTasks(s services.TaskServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id   string
			resp []byte
			err  error
			task = entities.Task{}
		)

		switch r.Method {
		case http.MethodPost:
			err = json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
				return
			}
			id, err = s.AddTask(task)
			resp, _ = json.Marshal(entities.Result{Id: id})
		case http.MethodGet:
			task, err = s.GetTask(r.FormValue("id"))
			resp, _ = json.Marshal(task)
		case http.MethodPut:
			err = json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
				return
			}
			err = s.EditTask(task)
			resp, _ = json.Marshal(task)
		case http.MethodDelete:
			err = s.DeleteTask(r.FormValue("id"))
			resp, _ = json.Marshal(task)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.Result{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(resp)
	}
}
