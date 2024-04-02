package main

import (
	"fmt"
	"main/app/controller"
	"main/app/shared/database"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func getPort() int {
	port := 7540
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}

	return port
}

func main() {
	database.InstallDb()
	webDir := "./web"

	r := chi.NewRouter()
	r.Mount("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", controller.NextDateReadGET)
	r.Post("/api/task", controller.TaskAddPOST)
	r.Get("/api/tasks", controller.TasksReadGET)

	err := http.ListenAndServe(fmt.Sprintf(":%d", getPort()), r)
	if err != nil {
		panic(err)
	}
}
