package main

import (
	"fmt"
	"go_final_project/internal/database"
	"go_final_project/internal/entities"
	"go_final_project/internal/handlers"
	"go_final_project/internal/services"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Overload(".env")
	if err != nil {
		log.Fatalln(err.Error())
	}

	entities.EnvMap, err = godotenv.Read(".env")
	if err != nil {
		log.Fatalln(err.Error())
	}

	dbName := entities.EnvMap["TODO_DBFILE"]
	handlers.Db, err = database.CheckAndGetDB(dbName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer handlers.Db.Close()

	mux := http.NewServeMux()

	uiDir := entities.EnvMap["TODO_UIDIR"]
	mux.Handle("/", http.FileServer(http.Dir(uiDir)))

	mux.HandleFunc("/api/signin", handlers.Authorization)
	mux.HandleFunc("/api/tasks", services.CheckJWT(handlers.GetTasks))
	mux.HandleFunc("/api/task", services.CheckJWT(handlers.UpdateTasks))
	mux.HandleFunc("/api/nextdate", handlers.GetNextDate)
	mux.HandleFunc("/api/task/done", services.CheckJWT(handlers.DoneTask))

	log.Println("Scheduler is running ...")
	port := fmt.Sprintf(":%s", entities.EnvMap["TODO_PORT"])
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("error when starting the server: %s\n", err.Error())
	}
}
