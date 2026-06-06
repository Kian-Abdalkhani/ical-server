package render

import (
	"fmt"
	"strings"
	"time"

	"example.com/ical/db"
)

// vtimezone produces a minimal VTIMEZONE component for the given IANA timezone name.
// It emits a single STANDARD sub-component based on the current UTC offset,
// which satisfies most calendar clients (Outlook, Apple Calendar, Google Calendar).
// A full historical RRULE expansion is not required for RFC 5545 conformance when
// calendar clients can resolve TZID via the IANA database.
func vtimezone(tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return ""
	}
	_, offset := time.Now().In(loc).Zone()
	hours := offset / 3600
	mins := (offset % 3600) / 60
	sign := "+"
	if offset < 0 {
		sign = "-"
		hours = -hours
		mins = -mins
	}
	offsetStr := fmt.Sprintf("%s%02d%02d", sign, hours, mins)

	var b strings.Builder
	fmt.Fprintf(&b, "BEGIN:VTIMEZONE\r\nTZID:%s\r\nBEGIN:STANDARD\r\nTZOFFSETFROM:%s\r\nTZOFFSETTO:%s\r\nDTSTART:19700101T000000\r\nEND:STANDARD\r\nEND:VTIMEZONE\r\n", tz, offsetStr, offsetStr)
	return b.String()
}

func RenderICS(events []db.Event) string {
	now := time.Now().UTC().Format("20060102T150405Z")
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//Family Calendar//EN\r\n")
	b.WriteString("CALSCALE:GREGORIAN\r\n")
	b.WriteString("METHOD:PUBLISH\r\n")

	// Emit one VTIMEZONE block per unique timezone used by timed (non-all-day) events.
	seen := make(map[string]bool)
	for _, e := range events {
		if !e.AllDay && e.Timezone != "" && !seen[e.Timezone] {
			if _, err := time.LoadLocation(e.Timezone); err == nil {
				b.WriteString(vtimezone(e.Timezone))
				seen[e.Timezone] = true
			}
		}
	}

	for _, e := range events {
		b.WriteString("BEGIN:VEVENT\r\n")
		fmt.Fprintf(&b, "UID:%s\r\n", e.UUID)
		fmt.Fprintf(&b, "DTSTAMP:%s\r\n", now)
		fmt.Fprintf(&b, "SUMMARY:%s\r\n", e.Summary)
		fmt.Fprintf(&b, "DESCRIPTION:%s\r\n", e.Description)

		if e.Location.Valid && e.Location.String != "" {
			fmt.Fprintf(&b, "LOCATION:%s\r\n", e.Location.String)
		}

		if e.AllDay {
			// RFC 5545 §3.6.1: all-day events use VALUE=DATE with just YYYYMMDD.
			fmt.Fprintf(&b, "DTSTART;VALUE=DATE:%s\r\n", e.Start.ICSDate())
			fmt.Fprintf(&b, "DTEND;VALUE=DATE:%s\r\n", e.End.ICSDate())
		} else if e.Timezone != "" {
			if _, err := time.LoadLocation(e.Timezone); err == nil {
				// RFC 5545 §3.6.1: timed event with known timezone uses TZID parameter.
				fmt.Fprintf(&b, "DTSTART;TZID=%s:%s\r\n", e.Timezone, e.Start.ICSInLocation(e.Timezone))
				fmt.Fprintf(&b, "DTEND;TZID=%s:%s\r\n", e.Timezone, e.End.ICSInLocation(e.Timezone))
			} else {
				// Invalid timezone stored: fall back to UTC.
				fmt.Fprintf(&b, "DTSTART:%s\r\n", e.Start.ICS())
				fmt.Fprintf(&b, "DTEND:%s\r\n", e.End.ICS())
			}
		} else {
			// No timezone: emit UTC.
			fmt.Fprintf(&b, "DTSTART:%s\r\n", e.Start.ICS())
			fmt.Fprintf(&b, "DTEND:%s\r\n", e.End.ICS())
		}

		b.WriteString("END:VEVENT\r\n")
	}
	b.WriteString("END:VCALENDAR")
	return b.String()
}
