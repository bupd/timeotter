package cron

import (
	"fmt"
	"log"
	"os/user"
	"runtime"

	"github.com/bupd/timeotter/pkg/utils"
)

func AddCrons(cronJob, cmdToExec string) error {
	cronLocation := GetCronLocation()

	// linux := "/var/spool/cron/bupd"
	// termux := "/data/data/com.termux/files/usr/var/spool/cron/u0_a323"

	// Shell script to backup the crontab and remove cron jobs below the comment
	script := fmt.Sprintf(`
    #!/bin/bash

    # Add a new cron job manually (This is your cron job command)
    echo "%s %s" >> "%s"
    `, cronJob, cmdToExec, cronLocation)

	err := utils.ExecuteShellCommand(script)
	if err != nil {
		log.Fatalf("failed to execute cron add command: %v", err)
	}

	return nil
}

func GetCronLocation() string {
	var cronLocation string

	curruser, _ := user.Current()

	if runtime.GOOS == "android" && runtime.GOARCH == "arm64" {
		termux := "/data/data/com.termux/files/usr/var/spool/cron"
		cronLocation = fmt.Sprintf("%s/%s", termux, curruser.Username)
	}
	if runtime.GOOS == "linux" {
		linux := "/var/spool/cron"
		cronLocation = fmt.Sprintf("%s/%s", linux, curruser.Username)
	}

	return cronLocation
}

func ClearCronJobs(backupFile string, cronMarker string) error {
	cronLocation := GetCronLocation()

	script := fmt.Sprintf(`
  #!/bin/bash

	# Backup current crontab
	crontab -l > "%s"

	# Remove all cron jobs below the comment
	awk '/%s/{f=1} !f' <(crontab -l) | crontab -
	echo "%s" >> "%s"
    `, backupFile, cronMarker, cronMarker, cronLocation)

	// Execute the shell script
	err := utils.ExecuteShellCommand(script)
	if err != nil {
		log.Fatalf("failed to execute cron remove command: %v", err)
	}

	return nil
}
