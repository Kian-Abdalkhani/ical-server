package event

import (
	"fmt"
	"time"

	"example.com/ical/db"
)

type Event struct {
	UUID        string     `json:"uuid" db:"uuid"`
	Summary     string     `json:"summary" db:"summary"`
	Location    string     `json:"location" db:"location"`
	Description string     `json:"description" db:"description"`
	Start       CustomTime `json:"start" db:"start"`
	End         CustomTime `json:"end" db:"end"`
	CreatedAt   time.Time  `db:"created_at"`
}

func NewEvent(uuid string, summmary string, location string, description string, start time.Time, end time.Time) *Event {
	return &Event{
		UUID:        uuid,
		Summary:     summmary,
		Location:    location,
		Description: description,
		Start:       CustomTime{start.UTC()},
		End:         CustomTime{end.UTC()},
		CreatedAt:   time.Now().UTC(),
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

	// if event already exists in the database
	_, err := GetEventByID(e.UUID)
	if err == nil {
		return fmt.Errorf("event id already exists in database")
	}

	insertSQL := `INSERT INTO events(
	uuid, summary, location, description, start, end, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = db.DB.Exec(
		insertSQL,
		e.UUID, e.Summary, e.Location,
		e.Description, e.Start, e.End,
		e.CreatedAt,
	)

	return err
}
