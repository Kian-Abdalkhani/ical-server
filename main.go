package main

import (
	"cmp"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"example.com/ical/db"
	"example.com/ical/handlers"
)

//go:embed schema.sql
var ddl string

func main() {
	ctx := context.Background()

	database, err := sql.Open("sqlite", ":memory:")
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
