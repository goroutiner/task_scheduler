package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go_final_project/internal/database"
	"io"
	"net/http"
	"strconv"
	"time"
)

// PostTask добавляет задачу с параметрами, полученными из тела запроса, в таблицу scheduler.
func PostTask(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		newTask Task
		id      string
		data    []byte
		err     error
	)

	body := r.Body
	defer body.Close()

	data, err = io.ReadAll(body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &newTask)
	if err != nil {
		return "", err
	}

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
			newTask.Date, err = NextDate(now, newTask.Date, newTask.Repeat)
			if err != nil {
				return "", err
			}
		}
	}

	task := &database.Task{
		Date:    newTask.Date,
		Title:   newTask.Title,
		Comment: newTask.Comment,
		Repeat:  newTask.Repeat,
	}

	id, err = task.PostTask(db)
	if err != nil {
		return "", err
	}

	return id, nil
}

// EditTask изменяет пармаетры задачи, полученные из тела запроса, в таблице scheduler.
func EditTask(db *sql.DB, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var (
		oldTask    Task
		err        error
		data, resp []byte
	)

	body := r.Body
	defer body.Close()

	data, err = io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &oldTask)
	if err != nil {
		return nil, err
	}

	if oldTask.Date == "" {
		oldTask.Date = time.Now().Format("20060102")
	}

	if oldTask.Title == "" {
		return nil, errors.New("the task title is empty")
	}

	if _, err := strconv.Atoi(oldTask.Id); err != nil {
		return nil, errors.New("the id is not specified or is specified not correctly")
	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	date, err := time.Parse("20060102", oldTask.Date)
	if err != nil {
		return nil, err
	}

	if date.Before(now) {
		if oldTask.Repeat == "" {
			oldTask.Date = now.Format("20060102")
		} else {
			oldTask.Date, err = NextDate(now, oldTask.Date, oldTask.Repeat)
			if err != nil {
				return nil, err
			}
		}
	}

	updateTask := &database.Task{
		Id:      oldTask.Id,
		Date:    oldTask.Date,
		Title:   oldTask.Title,
		Comment: oldTask.Comment,
		Repeat:  oldTask.Repeat,
	}

	err = updateTask.UpdateTask(db)
	if err != nil {
		return nil, err
	}

	resp, _ = json.Marshal(Task{})
	return resp, nil
}

// DeleteTask удаляет задачу с id, полученным из параметров запроса (...?id=...).
func DeleteTask(db *sql.DB, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var (
		id   string
		resp []byte
		err  error
	)

	if r.FormValue("id") == "" {
		return nil, errors.New("the id is not specified or is specified not correctly")
	}

	id = r.FormValue("id")
	err = database.DeleteTask(db, id)
	if err != nil {
		return nil, err
	}

	resp, _ = json.Marshal(Task{})
	return resp, nil
}
