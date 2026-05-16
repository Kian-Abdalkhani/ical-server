package db

import (
	"database/sql"
	"fmt"
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

	DB, err := sql.Open("sqlite3", DB_PATH)

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
			location TEXT
			description TEXT
			start DATETIME NOT NULL
			end DATETIME NOT NULL
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table")
	}

}
