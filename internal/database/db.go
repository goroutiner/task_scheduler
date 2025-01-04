package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksStore struct {
	Db *sql.DB
}

// CreateDB создает таблицу scheduler.
func CreateDB(db *sql.DB) {
	_, _ = db.Exec(`CREATE TABLE scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date INTEGER NOT NULL DEFAULT 0,
        title TEXT NOT NULL DEFAULT "",
        comment TEXT NOT NULL DEFAULT "",
        repeat VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX scheduler_date ON scheduler (date);
		`)
}

// CheckDB проверяет наличие БД с указанным именем, если нет, то создастся.
func CheckDB(fileName string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}
	dbFile := filepath.Join(path, "scheduler.db")

	var has bool
	if _, err := os.Stat(dbFile); err == nil {
		has = true
	}

	if !has {
		log.Println("Creating a database ...")
		file, err := os.Create(dbFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer file.Close()
	}
}

// Метод PostTask добавляет задачу, с указанными параметрами в таблицу scheduler.
func (tasks *Task) PostTask(db *sql.DB) (string, error) {
	row, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", tasks.Date), sql.Named("title", tasks.Title),
		sql.Named("comment", tasks.Comment), sql.Named("repeat", tasks.Repeat))
	if err != nil {
		return "", err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprint(id), nil
}

// GetTasks получат все задач из таблицы scheduler.
func GetTasks(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT * FROM scheduler ORDER BY date`)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// SearchTasks находит все задачи, содержащие строку или подстроку target
// в полях title или comment таблицы scheduler, если в url запроса есть параметр
// search (...?search=...). Если парметр не указан, то получаем все имеющиеся задачи.
func SearchTasks(db *sql.DB, target string) (*sql.Rows, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if date, err := time.Parse("02.01.2006", target); err == nil {
		dateInFormat := date.Format("20060102")

		rows, err = db.Query("SELECT * FROM scheduler WHERE date = ?", dateInFormat)
		if err != nil {
			return nil, err
		}

		return rows, nil
	}

	target = fmt.Sprint("%" + target + "%")
	rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", target))
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// SearchTask получает задачу с переданным id из таблицы scheduler.
func SearchTask(db *sql.DB, id string) *sql.Row {
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)

	return row
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
