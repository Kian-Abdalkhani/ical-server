package main

import (
	"cmp"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"example.com/ical/db"
	"example.com/ical/handlers"
)

//go:embed schema.sql
var ddl string

func main() {
	ctx := context.Background()

	fn := filepath.Join(".", "data", "ical.db")
	database, err := sql.Open("sqlite", fn)
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}

	if _, err := database.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries := db.New(database)

	port := cmp.Or(os.Getenv("PORT"), "8080")
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	handlers.Handle(queries)
	fmt.Printf("Listening on %s\n", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
