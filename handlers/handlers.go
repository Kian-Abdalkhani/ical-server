package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/ical/event"
	"example.com/ical/render"
)

func Handle() {

	// Serves static frontend files
	http.HandleFunc("/", webUIHandler)

	http.HandleFunc("/api/events", eventHandler)

	http.HandleFunc("/family.ics", calendarHandler)
}

func calendarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=family.ics")

	var events_list []event.Event
	// events := getEventsFromDB() //pull events from DB
	events_list = append(events_list)
	fmt.Fprint(w, render.RenderICS(events_list))
}

func webUIHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	var created_event event.Event
	if err := json.NewDecoder(r.Body).Decode(&created_event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Received event: %+v\n", created_event)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created_event)
}
