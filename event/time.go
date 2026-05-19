package event

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// CustomTime wraps time.Time with JSON, SQLite, and iCal serialization support.
// All times are normalized to UTC on parse.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON accepts RFC3339 and the "2006-01-02T15:04" format produced
// by HTML <input type="datetime-local">.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04"} {
		if t, err := time.Parse(layout, s); err == nil {
			ct.Time = t.UTC()
			return nil
		}
	}
	return fmt.Errorf("cannot parse time %q", s)
}

// MarshalJSON outputs UTC RFC3339.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.UTC().Format(time.RFC3339) + `"`), nil
}

// Value implements driver.Valuer, storing time as a UTC RFC3339 string in SQLite.
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.UTC().Format(time.RFC3339), nil
}

// Scan implements sql.Scanner, reading a UTC time string from SQLite.
func (ct *CustomTime) Scan(src any) error {
	switch v := src.(type) {
	case string:
		for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05"} {
			if t, err := time.Parse(layout, v); err == nil {
				ct.Time = t.UTC()
				return nil
			}
		}
		return fmt.Errorf("cannot scan time %q", v)
	case []byte:
		return ct.Scan(string(v))
	case int64:
		ct.Time = time.Unix(v, 0).UTC()
		return nil
	case nil:
		ct.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("unsupported time source type: %T", src)
	}
}

// ICS returns the time formatted per RFC 5545 (UTC, basic format).
func (ct CustomTime) ICS() string {
	return ct.UTC().Format("20060102T150405Z")
}
