// Package main is the entry point for the timeotter CLI application.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	cal "github.com/bupd/timeotter/pkg/calendar"
	"github.com/bupd/timeotter/pkg/config"
	"github.com/bupd/timeotter/pkg/oauth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	calendarID           string
	cmdToExec            string
	maxRes               int64
	tokenFile            string
	credentialsFile      string
	backupFile           string
	triggerBeforeMinutes int
	cronMarker           string
	showDeleted          bool
)

func main() {
	conf := config.GetConfig()

	calendarID = conf.CalendarID
	cmdToExec = conf.CmdToExec
	maxRes = conf.MaxRes
	tokenFile = conf.TokenFile
	credentialsFile = conf.CredentialsFile
	backupFile = conf.BackupFile
	triggerBeforeMinutes = conf.TriggerBeforeMinutes
	cronMarker = conf.CronMarker
	showDeleted = conf.ShowDeleted

	ctx := context.Background()
	b, err := os.ReadFile(filepath.Clean(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := oauth.GetClient(config, tokenFile)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// calList := srv.CalendarList.List()
	// kumar := srv.CalendarList.List().Fields()
	// Marshal the struct to a JSON string
	// jsonData, err := json.MarshalIndent(calList, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshaling struct: %v", err)
	// }

	// // Print the marshaled output
	// fmt.Println(string(jsonData))

	// jsonDatas, err := json.MarshalIndent(kumar, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshaling struct: %v", err)
	// }

	// Print the marshaled output
	// fmt.Println(string(jsonDatas))
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(calendarID).ShowDeleted(showDeleted).
		SingleEvents(true).TimeMin(t).MaxResults(maxRes).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		cal.EventParser(events, cmdToExec, backupFile, cronMarker, triggerBeforeMinutes)
	}
}
