package render

import (
	"fmt"
	"strings"
	"time"

	"example.com/ical/event"
)

func RenderICS(events []event.Event) string {
	now := time.Now().UTC().Format("20060102T150405Z")
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//Family Calendar//EN\r\n")
	b.WriteString("CALSCALE:GREGORIAN\r\n")
	b.WriteString("METHOD:PUBLISH\r\n")

	for _, e := range events {
		b.WriteString("BEGIN:VEVENT\r\n")
		fmt.Fprintf(&b,
			"UID:%s\r\nDTSTAMP:%s\r\nSUMMARY:%s\r\nDESCRIPTION:%s\r\nDTSTART:%s\r\nDTEND:%s\r\n",
			e.UID, now,
			e.Summary, e.Description,
			e.Start.Time.UTC().Format("20060102T150405Z"),
			e.End.Time.UTC().Format("20060102T150405Z"),
		)
		b.WriteString("END:VEVENT\r\n")
	}
	b.WriteString("END:VCALENDAR")
	return b.String()
}
