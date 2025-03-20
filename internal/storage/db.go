package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"task_scheduler/internal/config"
	"task_scheduler/internal/entities"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sqlx.DB
}

// NewDatabaseConection возвращает структуру соединения с БД.
func NewDatabaseConection(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// NewSqliteStore создает таблицу "scheduler" в БД для хранения посылок в режиме "sqlite".
func NewSqliteStore(fileName string) (*sqlx.DB, error) {
	var (
		db *sqlx.DB
	)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Println("Creating a database ...")
		_, err := os.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
	}

	db, err := sqlx.Open("sqlite", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date INTEGER NOT NULL DEFAULT 0,
        title TEXT NOT NULL DEFAULT '',
        comment TEXT NOT NULL DEFAULT '',
        repeat VARCHAR(128) NOT NULL DEFAULT ''
    );

	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);
		`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, err
}

// NewPostgresStore создает таблицу "scheduler" в БД для хранения посылок в режиме "postgres".
func NewPostgresStore(psqlUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", psqlUrl)
	if err != nil {
		return nil, err
	}

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS scheduler (
        id SERIAL PRIMARY KEY,
        date INTEGER NOT NULL DEFAULT 0,
        title TEXT NOT NULL DEFAULT '',
        comment TEXT NOT NULL DEFAULT '',
        repeat VARCHAR(128) NOT NULL DEFAULT ''
    );

	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);
		`)

	return db, err
}

// Метод PostTask добавляет задачу с указанными параметрами в таблицу scheduler.
func (s *Storage) PostTask(task entities.Task) (string, error) {
	var id int

	if config.Mode == "postgres" {
		query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES ($1, $2, $3, $4)
			  RETURNING id;`
		row := s.db.QueryRow(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err := row.Err(); err != nil {
			return "", err
		}

		row.Scan(&id)
	} else {
		query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES (:date, :title, :comment, :repeat)`

		_, err := s.db.NamedExec(query, task)
		if err != nil {
			return "", fmt.Errorf("failed to insert task: %w", err)
		}

		err = s.db.Get(&id, "SELECT last_insert_rowid()")
		if err != nil {
			return "", fmt.Errorf("failed to get last insert ID: %w", err)
		}
	}

	return fmt.Sprint(id), nil
}

// GetTasks получат все существующие задач из таблицы scheduler.
func (s *Storage) GetTasks() ([]entities.Task, error) {
	tasks := []entities.Task{}
	query := `SELECT * FROM scheduler ORDER BY date`

	err := s.db.Select(&tasks, query)

	return tasks, err
}

// SearchTasks возвращает все задачи, содержащие строку или подстроку target
// в полях title или comment из таблицы scheduler.
func (s *Storage) SearchTasks(target string) ([]entities.Task, error) {
	var (
		tasks = []entities.Task{}
		query string
	)

	if date, err := time.Parse("02.01.2006", target); err == nil {
		dateInFormat := date.Format("20060102")

		if config.Mode == "postgres" {
			query = `SELECT * FROM scheduler WHERE date = $1`
		} else {
			query = `SELECT * FROM scheduler WHERE date = ?`
		}

		err := s.db.Select(&tasks, query, dateInFormat)

		return tasks, err
	}

	target = fmt.Sprint("%" + target + "%")

	if config.Mode == "postgres" {
		query = `SELECT * FROM scheduler WHERE title ILIKE $1 OR comment ILIKE $1 ORDER BY date`
	} else {
		query = `SELECT * FROM scheduler WHERE LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?) ORDER BY date`
	}

	err := s.db.Select(&tasks, query, target, target)

	return tasks, err
}

// SearchTask получает задачу по id из таблицы scheduler.
func (s *Storage) SearchTask(id string) (entities.Task, error) {
	var (
		task  = entities.Task{}
		query string
	)

	if config.Mode == "postgres" {
		query = `SELECT * FROM scheduler WHERE id = $1`
	} else {
		query = `SELECT * FROM scheduler WHERE id = ?`
	}

	err := s.db.Get(&task, query, id)

	return task, err
}

// Метод UpdateTask обновляет задачу по переданным параметрам в таблице scheduler.
func (s *Storage) UpdateTask(task entities.Task) error {
	var (
		exists            bool
		queryCheck, query string
	)

	if config.Mode == "postgres" {
		queryCheck = `SELECT EXISTS (SELECT 1 FROM scheduler WHERE id = $1)`

		query = `UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5`
	} else {
		queryCheck = `SELECT EXISTS (SELECT 1 FROM scheduler WHERE id = ?)`

		query = `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	}

	if s.db.Get(&exists, queryCheck, task.Id); !exists {
		return errors.New("there is no task with the specified id")
	}

	_, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)

	return err
}

// DeleteTask удаляет задачу по id из таблицы scheduler.
func (s *Storage) DeleteTask(id string) error {
	var (
		exists            bool
		queryCheck, query string
	)

	if config.Mode == "postgres" {
		queryCheck = `SELECT EXISTS (SELECT 1 FROM scheduler WHERE id = $1)`

		query = `DELETE FROM scheduler WHERE id = $1`
	} else {
		queryCheck = `SELECT EXISTS (SELECT 1 FROM scheduler WHERE id = ?)`

		query = `DELETE FROM scheduler WHERE id = ?`
	}

	if s.db.Get(&exists, queryCheck, id); !exists {
		return errors.New("there is no task with the specified id")
	}

	_, err := s.db.Exec(query, id)

	return err
}
