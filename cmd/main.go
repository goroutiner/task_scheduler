package main

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/database"
	"go_final_project/internal/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const webDir = "../ui"

func main() {
	err := godotenv.Overload("../.env")
	if err != nil {
		log.Println(err.Error())
	}

	database.CheckDB("scheduler.db")

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	database.CreateDB(db)

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	mux.HandleFunc("/api/signin", handlers.Authorization)
	mux.HandleFunc("/api/tasks", handlers.GetTasks)
	mux.HandleFunc("/api/task", handlers.UpdateTasks)
	mux.HandleFunc("/api/nextdate", handlers.GetNextDate)
	mux.HandleFunc("/api/task/done", handlers.DoneTask)

	envMap, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatalln(err.Error())
	}
	port := fmt.Sprintf(":%s", envMap["TODO_PORT"]) 

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %s\n", err.Error())
	}
}
