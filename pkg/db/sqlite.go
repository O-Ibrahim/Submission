package db

import (
	"database/sql"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func GetDBDriver() (*sql.DB, error) {
	log.Println("getting db driver")
	d, err := sql.Open("sqlite3", "file:demo.db")
	if err != nil {
		return nil, err
	}
	log.Println("pinging db")
	if err := d.Ping(); err != nil {
		return nil, err
	}

	if err := CreateSchema(d); err != nil {
		return nil, err
	}
	return d, nil
}

func CreateSchema(db *sql.DB) error {
	log.Println("creating schema")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		command TEXT,
		args TEXT,
		status TEXT,
		logfile TEXT,
		created_at INTEGER,
		updated_at INTEGER
	)`)
	return err
}
