package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"example.com/ical/db"
	"example.com/ical/handlers"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)

	handlers.Handle()
	// Start LAN server on the specified port
	fmt.Printf("Listening on port %d\n", port)
	http.ListenAndServe(addr, nil)
}
