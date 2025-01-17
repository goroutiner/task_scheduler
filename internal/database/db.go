package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/internal/entities"
	"log"
	"os"
	"time"
)

type Task entities.Task

// CheckAndGetDB проверяет наличие БД с указанным именем, если нет, то создастся.
func CheckAndGetDB(fileName string) (*sql.DB, error) {
	var (
		has bool
		db  *sql.DB
	)

	if _, err := os.Stat(fileName); err == nil {
		has = true
	}

	if !has {
		log.Println("Creating a database ...")
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		db, err := sql.Open("sqlite", fileName)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`CREATE TABLE scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date INTEGER NOT NULL DEFAULT 0,
        title TEXT NOT NULL DEFAULT "",
        comment TEXT NOT NULL DEFAULT "",
        repeat VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX scheduler_date ON scheduler (date);
		`)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", fileName)
	return db, err
}

// Метод PostTask добавляет задачу, с указанными параметрами в таблицу scheduler.
func (tasks *Task) PostTask(db *sql.DB) (string, error) {
	var id  int64

	row, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", tasks.Date), sql.Named("title", tasks.Title),
		sql.Named("comment", tasks.Comment), sql.Named("repeat", tasks.Repeat))
	if err != nil {
		return "", err
	}

	id, err = row.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error while getting last id: %w", err)
	}

	return fmt.Sprint(id), nil
}

// GetTasks получат все задач из таблицы scheduler.
func GetTasks(db *sql.DB) ([]entities.Task, error) {
	var (
		task  entities.Task
		tasks = []entities.Task{}
	)

	rows, err := db.Query(`SELECT * FROM scheduler ORDER BY date`)
	if err != nil {
		return nil, err
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

// SearchTasks находит все задачи, содержащие строку или подстроку target
// в полях title или comment таблицы scheduler, если в url запроса есть параметр
// search (...?search=...). Если парметр не указан, то получаем все имеющиеся задачи.
func SearchTasks(db *sql.DB, target string) ([]entities.Task, error) {
	var (
		task  entities.Task
		tasks = []entities.Task{}
	)

	if date, err := time.Parse("02.01.2006", target); err == nil {
		dateInFormat := date.Format("20060102")

		rows, err := db.Query("SELECT * FROM scheduler WHERE date = ?", dateInFormat)
		if err != nil {
			return nil, err
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

	target = fmt.Sprint("%" + target + "%")
	rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", target))
	if err != nil {
		return nil, err
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

// SearchTask получает задачу с переданным id из таблицы scheduler.
func SearchTask(db *sql.DB, id string) (entities.Task, error) {
	var task entities.Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return entities.Task{}, err
	}

	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return entities.Task{}, err
	}

	return task, nil
}

// Метод UpdateTask обновляет задачу по переданным параметрам из таблицы scheduler.
func (tasks *Task) UpdateTask(db *sql.DB) error {
	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", tasks.Date),
		sql.Named("title", tasks.Title),
		sql.Named("comment", tasks.Comment),
		sql.Named("repeat", tasks.Repeat),
		sql.Named("id", tasks.Id))
	if err != nil {
		return err
	}

	num, err := res.RowsAffected()
	if num == 0 {
		err = errors.New("there is no task with the specified id")
	}

	return err
}

// DeleteTask удаляет задачу с переданным id из таблицы scheduler.
func DeleteTask(db *sql.DB, id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	num, err := res.RowsAffected()
	if num == 0 {
		err = errors.New("there is no task with the specified id")
	}

	return err
}
