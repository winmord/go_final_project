package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	envDbFilePath := os.Getenv(" TODO_DBFILE ")
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

func main() {
	installDb()
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(fmt.Sprintf(":%d", getPort()), nil)
	if err != nil {
		panic(err)
	}
}
