package services

import (
	"database/sql"
	"errors"
	"go_final_project/internal/database"
	"go_final_project/internal/entities"
	"net/http"
)

// GetTasks получает задачи, содержащие строку или подсторку,
// полученную из параметров запроса.
func GetTasks(db *sql.DB, r *http.Request) ([]entities.Task, error) {
	var (
		target string
		tasks  = []entities.Task{}
		err    error
	)

	// Поиск записи по значению параметра search, если он указан,
	// в противном случае показываются все задачи.
	if target = r.FormValue("search"); target != "" {
		tasks, err = database.SearchTasks(db, target)
		if err != nil {
			return nil, err
		}
	} else {
		tasks, err = database.GetTasks(db)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

// GetTask получает задачу по id, полученного из параметров запроса.
func GetTask(db *sql.DB, w http.ResponseWriter, r *http.Request) (entities.Task, error) {
	var (
		task entities.Task
		id   string
		err  error
	)

	if r.FormValue("id") == "" {
		return entities.Task{}, errors.New("the id is not specified or is specified not correctly")
	}

	id = r.FormValue("id")
	task, err = database.SearchTask(db, id)
	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}
