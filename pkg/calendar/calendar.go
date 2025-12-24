package calendar

import (
	"fmt"
	"log"
	"time"

	"github.com/bupd/timeotter/pkg/cron"
	"google.golang.org/api/calendar/v3"
)

// for parsing calendar events
func EventParser(events *calendar.Events, cmdToExec string, backupFile string, cronMarker string, triggerBeforeMinutes int) {
	err := cron.ClearCronJobs(backupFile, cronMarker)
	if err != nil {
		log.Fatalf("clearing cron jobs failed: %v", err)
	}

	for _, item := range events.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}
		// fmt.Printf("cron string: %s:- ", ConvertTimeToCron(date))
		// fmt.Printf("%v (%v)\n", item.Summary, date)
		cronStr := ConvertTimeToCron(date, triggerBeforeMinutes)
		err := cron.AddCrons(cronStr, cmdToExec)
		if err != nil {
			log.Fatalf("unable to add crons: %v", err)
		}
	}
}

// ConvertTimeToCron takes a time in string format and returns the cron expression.
// triggerBeforeMinutes specifies how many minutes before the event to trigger.
func ConvertTimeToCron(timeStr string, triggerBeforeMinutes int) string {
	// Parse the input time string
	// Try multiple layouts to parse the time string
	layouts := []string{
		"2006-01-02T15:04:05-07:00", // e.g., 2025-02-02T20:29:00+05:30
		"2006/01/02T15:04:05-07:00", // e.g., 2025/02/02T20:29:00+05:30 (slashes in date)
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			break
		}
	}

	// Subtract triggerBeforeMinutes from the given time
	t = t.Add(-time.Duration(triggerBeforeMinutes) * time.Minute)

	// Extract time components
	_ = t.Second()
	mins := t.Minute()
	hour := t.Hour()
	day := t.Day()
	month := int(t.Month())
	weekday := int(t.Weekday())

	// Construct the cron string (without seconds as it is optional)
	cron := fmt.Sprintf("%d %d %d %d %d", mins, hour, day, month, weekday)
	return cron
}
