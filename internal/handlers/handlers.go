package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/internal/database"
	"go_final_project/internal/entities"
	"go_final_project/internal/services"
	"io"
	"log"
	"net/http"
	"time"
)

var Db *sql.DB

// Authorization получает пароль через https, и если он является валидным,
// то отправляет http ответ, содержащий сгенерированный токен.
func Authorization(w http.ResponseWriter, r *http.Request) {
	var (
		signedToken string
		resp        []byte
		// В переменной takenMap будет содержаться структура с паролем, полученная через https.
		takenMap map[string]string
	)

	body := r.Body
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	err = json.Unmarshal(data, &takenMap)
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	signedToken, err = services.GetJWT(takenMap)
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}
	resp, _ = json.Marshal(entities.Result{Token: signedToken})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// GetNextDate получает занчения параметров now, date, repeat из парметров запроса и
// с их помощью возвращает http ответ, содержащий следующую ближайшую дату.
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	res, err := services.NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// GetTasks возвращает http ответ, содержащий список всех сущетсвующих задач.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		tasks = []entities.Task{}
		resp  []byte
	)

	tasksTmp, err := services.GetTasks(Db, r)
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	for _, v := range tasksTmp {
		tasks = append(tasks, entities.Task{Id: v.Id, Date: v.Date, Title: v.Title, Comment: v.Comment, Repeat: v.Repeat})
	}
	resp, _ = json.Marshal(entities.Result{Tasks: tasks})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// DoneTask завершает или обноввляет дату задачи, если поле repeat не пустое и
// возвращает пустой json http ответа в случае успешного завершения.
func DoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	var (
		task entities.Task
		resp []byte
		id   string
		err  error
	)

	if r.FormValue("id") == "" {
		resp, _ = json.Marshal(entities.Result{Error: "id не указан или указан некорректно"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		log.Println("Id не указан или указан некорректно")
		return
	}

	id = r.FormValue("id")
	task, err = database.SearchTask(Db, id)
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	if task.Repeat != "" {
		nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			resp, _ = json.Marshal(entities.Result{Error: err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			log.Println(err.Error())
			return
		}

		updateTask := &database.Task{
			Id:      task.Id,
			Date:    nextDate,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}

		err = updateTask.UpdateTask(Db)
		if err != nil {
			resp, _ = json.Marshal(entities.Result{Error: err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			log.Println(err.Error())
			return
		}
	} else {
		err = database.DeleteTask(Db, id)
		if err != nil {
			resp, _ = json.Marshal(entities.Result{Error: err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			log.Println(err.Error())
			return
		}
	}

	resp, _ = json.Marshal(entities.Task{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// UpdateTasks возвращает id новой задачи в случае POST запроса,
// вовзращает задачу в случае GET запроса
// и возвращает пустой json http ответа в случае успешного завершения изменения задачи
// с помощью метода PUT и DELETE.
func UpdateTasks(w http.ResponseWriter, r *http.Request) {
	var (
		id   string
		resp []byte
		err  error
		task entities.Task
	)

	switch r.Method {
	case http.MethodPost:
		id, err = services.PostTask(Db, w, r)
		resp, _ = json.Marshal(entities.Result{Id: id})
	case http.MethodGet:
		task, err = services.GetTask(Db, w, r)
		resp, _ = json.Marshal(entities.Task{Id: task.Id, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat})
	case http.MethodPut:
		resp, err = services.EditTask(Db, w, r)
	case http.MethodDelete:
		resp, err = services.DeleteTask(Db, w, r)
	}
	if err != nil {
		resp, _ = json.Marshal(entities.Result{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}
