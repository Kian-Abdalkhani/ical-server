# Copilot Instructions

## Project Overview

`ical-server` is a Go HTTP server that manages calendar events stored in SQLite and exposes them as an iCal (`.ics`) feed consumable by calendar apps. It also serves a minimal web UI for creating events.

## Build & Run

```bash
# Build
go build ./...

# Run (defaults to port 8080)
go run main.go

# Run on a custom port
PORT=3000 go run main.go

# Run tests
go test ./...

# Run tests for a single package
go test example.com/ical/event
```

## Architecture

```
main.go           → initializes DB, registers routes, starts HTTP server
db/db.go          → opens the SQLite database at ./data, creates the events table
event/event.go    → Event struct + CRUD helpers (Save, Update, GetEventByID, GetAllEvents)
handlers/handlers.go → HTTP route registration and request handling
render/render.go  → renders []event.Event to RFC 5545 iCal format
static/index.html → single-file web UI (vanilla JS, posts JSON to /api/events)
```

### HTTP routes

| Route | Method | Description |
|---|---|---|
| `/` | GET | Serves `static/index.html` |
| `/api/events` | POST | Accepts a JSON `Event` body |
| `/family.ics` | GET | Returns the iCal feed (`Content-Type: text/calendar`) |

### Data flow

Event creation: `static/index.html` → `POST /api/events` → `handlers.eventHandler` → `event.Save()`

iCal export: `GET /family.ics` → `handlers.calendarHandler` → `event.GetAllEvents()` → `render.RenderICS()`

## Key Conventions

- **Module path**: `example.com/ical` (see `go.mod`). Use this for all internal imports.
- **SQLite driver**: `modernc.org/sqlite` (pure Go, no CGo). The global connection is `db.DB`.
- **iCal output**: Must use CRLF (`\r\n`) line endings per RFC 5545. All times must be UTC formatted as `20060102T150405Z`.
- **Time parsing**: `render.CustomTime` wraps `time.Time` with a custom `UnmarshalJSON` that accepts both `time.RFC3339` and `"2006-01-02T15:04"` (the format sent by `<input type="datetime-local">`).
- **Event identity**: Events use UUID strings as the primary key (`uuid TEXT PRIMARY KEY`). Callers are responsible for generating UUIDs (use `github.com/google/uuid`).
- **Database location**: The SQLite file is created at `./data` relative to the working directory.
