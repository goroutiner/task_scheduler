package services

import (
	"database/sql"
	"errors"
	"go_final_project/internal/database"
	"net/http"
)

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// GetTasks получает задачи, содержащие строку или подсторку, 
// полученную из параметров запроса.
func GetTasks(db *sql.DB, r *http.Request) ([]Task, error) {
	var (
		target string
		rows   *sql.Rows
		task   Task
		tasks  = []Task{}
		err    error
	)

	// Поиск записи по значению параметра search, если он указан,
	// в противном случае показываются все задачи.
	if target = r.FormValue("search"); target != "" {
		rows, err = database.SearchTasks(db, target)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = database.GetTasks(db)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает задачу по id, полученного из параметров запроса.
func GetTask(db *sql.DB, w http.ResponseWriter, r *http.Request) (Task, error) {
	var (
		task Task
		id   string
		row  *sql.Row
		err  error
	)

	if r.FormValue("id") == "" {
		return Task{}, errors.New("the id is not specified or is specified not correctly")
	}

	id = r.FormValue("id")
	row = database.SearchTask(db, id)
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
