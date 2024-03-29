package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Sheduler struct {
	Id      int
	Date    string
	Title   string
	Comment string
	Repeat  string
}

func getDbFilePath() string {
	//appPath, err := os.Executable()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//dbFilePath := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	dbFilePath := "scheduler.db"

	envDbFilePath := os.Getenv("TODO_DBFILE")
	if len(envDbFilePath) > 0 {
		dbFilePath = envDbFilePath
	}

	return dbFilePath
}

func createDbFile(dbFilePath string) (*sql.DB, error) {
	_, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE `scheduler` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` VARCHAR(8) NOT NULL, `title` VARCHAR(64) NOT NULL, `comment` VARCHAR(255) NOT NULL, `repeat` VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

func installDb() {
	dbFilePath := getDbFilePath()
	_, err := os.Stat(dbFilePath)

	var install bool
	if err != nil {
		install = true
	}

	if !install {
		dbFile, err := createDbFile(dbFilePath)
		if err != nil {
			log.Fatal(err)
		}
		createTable(dbFile)
	}
}

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

const datePattern string = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty string")
	} else if strings.Contains(repeat, "d ") {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", errors.New("maximum days count must be 400")
		}

		parsedDate, err := time.Parse(datePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format(datePattern), nil
	} else if repeat == "y" {
		parsedDate, err := time.Parse(datePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(datePattern), nil
	} else {
		return "", errors.New("repeat wrong format")
	}

}

func getNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(datePattern, r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := NextDate(now, date, repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error: %w", err).Error(), http.StatusBadRequest)
	}
}

func main() {
	installDb()
	webDir := "./web"

	r := chi.NewRouter()
	r.Mount("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", getNextDate)

	err := http.ListenAndServe(fmt.Sprintf(":%d", getPort()), r)
	if err != nil {
		panic(err)
	}
}
