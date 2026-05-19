package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

var DB *sql.DB

var DB_PATH = filepath.Join(".", "data")

func InitDB() {
	// Create a directory if necessary
	dir := DB_PATH
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("Could not create directory in path: %v", err))
	}

	var err error
	DB, err = sql.Open("sqlite", filepath.Join(DB_PATH, "ical.db"))

	if err != nil {
		panic("Could not connect to database")
	}

	// set limit on concurrent DB connections
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTable()
}

func createTable() {
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
			uuid TEXT PRIMARY KEY,
			summary TEXT NOT NULL,
			location TEXT,
			description TEXT,
			start TEXT NOT NULL,
			end TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
	);`

	_, err := DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table")
	}

}
