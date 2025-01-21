package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"os/user"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calList := srv.CalendarList.List()
	kumar := srv.CalendarList.List().Fields()
	// Marshal the struct to a JSON string
	jsonData, err := json.MarshalIndent(calList, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	// Print the marshaled output
	fmt.Println(string(jsonData))

	jsonDatas, err := json.MarshalIndent(kumar, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	// Print the marshaled output
	fmt.Println(string(jsonDatas))
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("bupdprasanth@gmail.com").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(2).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		eventParser(events)
	}
}

func eventParser(events *calendar.Events) {
	ClearCronJobs()

	for _, item := range events.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}
		fmt.Printf("cron string: %s:- ", ConvertTimeToCron(date))
		fmt.Printf("%v (%v)\n", item.Summary, date)
		cronStr := ConvertTimeToCron(date)
		err := AddCronJob(cronStr)
		if err != nil {
			// log.Fatal(err)
		} else {
			fmt.Println("Cron job added successfully!")
		}
	}
}

// ConvertTimeToCron takes a time in string format and returns the cron expression.
func ConvertTimeToCron(timeStr string) string {
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

	// Extract time components
	_ = t.Second()
	mins := t.Minute()
	hour := t.Hour()
	day := t.Day()
	month := int(t.Month())
	weekday := int(t.Weekday())

	// Construct the cron string (without seconds as it is optional)
	cron := fmt.Sprintf("%d %d %d %d %d *", mins, hour, day, month, weekday)
	return cron
}

// AddCronJob adds the cron job to the crontab.
func AddCronJob(cronString string) error {
	curruser, _ := user.Current()
	fmt.Println(curruser.Username)

	command := "mpv ~/video.mp4"
	// Build the cron job command
	cronJob := fmt.Sprintf("%s %s", cronString, command)

	// Use crontab command to add the cron job
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l; echo \"%s\") | crontab -", cronJob))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	return nil
}

func ClearCronJobs() error {
	// linux := "/var/spool/cron/bupd"
	// termux := "/data/data/com.termux/files/usr/var/spool/cron/u0_a323"

	// Shell script to backup the crontab and remove cron jobs below the comment
	script := `
	#!/bin/bash

	# Backup current crontab
	crontab -l > crontab_backup.txt

	# Remove all cron jobs below the comment
	awk '/# custom crons below this can be deleted/{f=1} !f' <(crontab -l) | crontab -
	echo "# custom crons below this can be deleted." >> "/var/spool/cron/bupd"
	`

	// should edit the above based on the linux or termux

	// Execute the shell script
	err := executeShellCommand(script)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Function to execute a shell command
func executeShellCommand(command string) error {
	// Create the command and pass it to a new shell
	cmd := exec.Command("bash", "-c", command)

	// Run the command and get any output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}

	// Print the output if the command was successful
	fmt.Println("Output:", string(output))
	return nil
}
