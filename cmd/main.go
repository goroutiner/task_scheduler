package main

import (
	"task_scheduler/internal/config"
	"task_scheduler/internal/entities"
	"task_scheduler/internal/handlers"
	"task_scheduler/internal/services"
	"task_scheduler/internal/storage"
	"time"

	"log"
	"net/http"
)

func main() {
	var store *storage.Storage

	switch config.Mode {
	case "sqlite":
		db, err := storage.NewSqliteStore(entities.DbFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer db.Close()

		store = storage.NewDatabaseConection(db)
		log.Println("Using SQLite storage")
	case "postgres":
		db, err := storage.NewPostgresStore(config.PsqlUrl)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer db.Close()

		store = storage.NewDatabaseConection(db)
		log.Println("Using PostgreSQL store")
	default:
		log.Fatalf("config.Mode is empty in /internal/config/setting.go")
	}

	taskService := services.GetTaskService(store)
	authService := services.GetAuthService()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(entities.UiDir)))

	mux.HandleFunc("GET /api/nextdate", handlers.GetNextDate(taskService))
	mux.HandleFunc("/api/task", services.CheckJWTMiddleware(handlers.UpdateTasks(taskService)))
	mux.HandleFunc("GET /api/tasks", services.CheckJWTMiddleware(handlers.GetTasks(taskService)))
	mux.HandleFunc("POST /api/task/done", services.CheckJWTMiddleware(handlers.DoneTask(taskService)))
	mux.HandleFunc("POST /api/signin", handlers.Authentication(authService))

	serv := &http.Server{
		Addr:         config.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("Scheduler is running ...")
	if err := serv.ListenAndServe(); err != nil {
		log.Fatalf("error when starting the server: %s\n", err.Error())
	}
}
