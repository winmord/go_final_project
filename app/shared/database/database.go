package database

import (
	"database/sql"
	"log"
	"main/app/model"
	"os"
)

var db *sql.DB

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

	db, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) {
	_, err := db.Exec(
		"CREATE TABLE `scheduler` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` VARCHAR(8) NULL, `title` VARCHAR(64) NOT NULL, `comment` VARCHAR(255) NULL, `repeat` VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

func InstallDb() {
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

func InsertTask(task model.Task) (int, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
