package event

import (
	"fmt"
	"strings"
	"time"

	"example.com/ical/db"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			ct.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse time %q", s)
}

type Event struct {
	UUID        string    `json:"uuid" db:"uuid"`
	Summary     string    `json:"summary" db:"summary"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	Start       time.Time `json:"start" db:"start"`
	End         time.Time `json:"end" db:"end"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewEvent(uuid string, summmary string, location string, description string, start time.Time, end time.Time) *Event {
	return &Event{
		UUID:        uuid,
		Summary:     summmary,
		Location:    location,
		Description: description,
		Start:       start,
		End:         end,
		CreatedAt:   time.Now(),
	}
}

func GetEventByID(uuid string) (*Event, error) {
	query := `SELECT uuid, summary, location, description, start, end, created_at
						FROM events WHERE uuid = ?`

	var e Event

	err := db.DB.QueryRow(query, uuid).Scan(
		&e.UUID, &e.Summary, &e.Location, &e.Description,
		&e.Start, &e.End, &e.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("event not found: %v", err)
	}

	return &e, nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT uuid, summary, location, description, start, end, created_at
						FROM events`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retreiving all events %v", err)
	}

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.UUID, &e.Summary, &e.Location, &e.Description,
			&e.Start, &e.End, &e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row %v", err)
		}

		events = append(events, e)
	}

	return events, nil
}

func (e *Event) Update() error {
	updateSQL := `UPDATE events SET summary = ?,
 location = ?, description = ?, start = ?, end = ?,
 created_at = ? WHERE uuid = ?`

	_, err := db.DB.Exec(
		updateSQL,
		e.Summary, e.Location, e.Description,
		e.Start, e.End, e.CreatedAt, e.UUID,
	)

	return err
}

func (e *Event) Save() error {
	insertSQL := `INSERT INTO events(
	uuid, summary, location, description, start, end, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(
		insertSQL,
		e.UUID, e.Summary, e.Location,
		e.Description, e.Start, e.End,
		e.CreatedAt,
	)

	return err
}
