package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cal "github.com/bupd/timeotter/pkg/calendar"
	"github.com/bupd/timeotter/pkg/config"
	"github.com/bupd/timeotter/pkg/oauth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	CalendarID string
	CmdToExec  string
	MaxRes     int64
	TokenFile  string
)

func main() {
	conf := config.GetConfig()

	CalendarID = conf.CalendarID
	CmdToExec = conf.CmdToExec
	MaxRes = conf.MaxRes
	TokenFile = conf.TokenFile

	dirname := config.GetHomeDir()
	credentialsFile := fmt.Sprintf("%s/.cal-credentials.json", dirname)
	TokenFile = fmt.Sprintf("%s/%s", dirname, TokenFile)

	ctx := context.Background()
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := oauth.GetClient(config, TokenFile)

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
	events, err := srv.Events.List(CalendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(MaxRes).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		cal.EventParser(events, CmdToExec)
	}
}
