package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"example.com/ical/components"
	"example.com/ical/db"
	"example.com/ical/render"
	"example.com/ical/timetype"
	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// Handle registers all HTTP routes, injecting the sqlc Queries into each handler.
func Handle(queries *db.Queries) {
	http.HandleFunc("GET /{$}", makePageHandler(queries))
	http.HandleFunc("GET /events", makeListHandler(queries))
	http.HandleFunc("POST /events", makeCreateHandler(queries))
	http.HandleFunc("GET /events/{uuid}/edit", makeEditFormHandler(queries))
	http.HandleFunc("GET /events/{uuid}/row", makeRowHandler(queries))
	http.HandleFunc("PUT /events/{uuid}", makeUpdateHandler(queries))
	http.HandleFunc("DELETE /events/{uuid}", makeDeleteHandler(queries))
	http.HandleFunc("GET /family.ics", makeCalendarHandler(queries))
}

func render_component(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// makePageHandler renders the full page with layout + event list + create form.
func makePageHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := queries.GetAllEvents(r.Context())
		if err != nil {
			http.Error(w, "failed to load events", http.StatusInternalServerError)
			return
		}
		render_component(w, r, components.Layout(events))
	}
}

// makeListHandler returns just the EventList fragment (used for HTMX partial refreshes).
func makeListHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := queries.GetAllEvents(r.Context())
		if err != nil {
			http.Error(w, "failed to load events", http.StatusInternalServerError)
			return
		}
		render_component(w, r, components.EventList(events))
	}
}

// makeCreateHandler handles form POSTs, creates a new event, returns an EventRow fragment.
func makeCreateHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}

		start, err := timetype.ParseCustomTime(r.FormValue("start"))
		if err != nil {
			http.Error(w, "invalid start time: "+err.Error(), http.StatusBadRequest)
			return
		}
		end, err := timetype.ParseCustomTime(r.FormValue("end"))
		if err != nil {
			http.Error(w, "invalid end time: "+err.Error(), http.StatusBadRequest)
			return
		}

		location := r.FormValue("location")
		var loc sql.NullString
		if location != "" {
			loc = sql.NullString{String: location, Valid: true}
		}

		params := db.CreateEventParams{
			UUID:        uuid.New().String(),
			Summary:     r.FormValue("summary"),
			Location:    loc,
			Description: r.FormValue("description"),
			Start:       start,
			End:         end,
			CreatedAt:   timetype.CustomTime{Time: time.Now().UTC()},
		}

		if err := queries.CreateEvent(r.Context(), params); err != nil {
			http.Error(w, "failed to create event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		event := db.Event{
			UUID:        params.UUID,
			Summary:     params.Summary,
			Location:    params.Location,
			Description: params.Description,
			Start:       params.Start,
			End:         params.End,
			CreatedAt:   params.CreatedAt,
		}
		render_component(w, r, components.EventRow(event))
	}
}

// makeEditFormHandler returns an EditForm fragment replacing the row.
func makeEditFormHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		event, err := queries.GetEventByID(r.Context(), id)
		if err != nil {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		render_component(w, r, components.EditForm(event))
	}
}

// makeRowHandler returns a single EventRow fragment (used by edit cancel).
func makeRowHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		event, err := queries.GetEventByID(r.Context(), id)
		if err != nil {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		render_component(w, r, components.EventRow(event))
	}
}

// makeUpdateHandler handles PUT form submissions, updates the event, returns updated EventRow.
func makeUpdateHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}

		start, err := timetype.ParseCustomTime(r.FormValue("start"))
		if err != nil {
			http.Error(w, "invalid start time: "+err.Error(), http.StatusBadRequest)
			return
		}
		end, err := timetype.ParseCustomTime(r.FormValue("end"))
		if err != nil {
			http.Error(w, "invalid end time: "+err.Error(), http.StatusBadRequest)
			return
		}

		location := r.FormValue("location")
		var loc sql.NullString
		if location != "" {
			loc = sql.NullString{String: location, Valid: true}
		}

		params := db.UpdateEventParams{
			Summary:     r.FormValue("summary"),
			Location:    loc,
			Description: r.FormValue("description"),
			Start:       start,
			End:         end,
			UUID:        id,
		}

		if err := queries.UpdateEvent(r.Context(), params); err != nil {
			http.Error(w, "failed to update event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		event, err := queries.GetEventByID(r.Context(), id)
		if err != nil {
			http.Error(w, "event not found after update", http.StatusInternalServerError)
			return
		}
		render_component(w, r, components.EventRow(event))
	}
}

// makeDeleteHandler deletes an event and returns an empty 200 (HTMX removes the row).
func makeDeleteHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		if err := queries.DeleteEvent(r.Context(), id); err != nil {
			http.Error(w, "failed to delete event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func makeCalendarHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=family.ics")

		events, err := queries.GetAllEvents(r.Context())
		if err != nil {
			http.Error(w, "failed to load events", http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprint(w, render.RenderICS(events))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
