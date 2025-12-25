# 🦦 TimeOtter 🦦

[![GitHub Release](https://img.shields.io/github/v/release/bupd/timeotter?style=flat-square&logo=github&color=blue)](https://github.com/bupd/timeotter/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bupd/timeotter?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/bupd/timeotter?style=flat-square&color=green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bupd/timeotter?style=flat-square)](https://goreportcard.com/report/github.com/bupd/timeotter)

[![Build](https://img.shields.io/github/actions/workflow/status/bupd/timeotter/build.yml?style=flat-square&logo=github&label=build)](https://github.com/bupd/timeotter/actions/workflows/build.yml)
[![Test](https://img.shields.io/github/actions/workflow/status/bupd/timeotter/test.yml?style=flat-square&logo=github&label=test)](https://github.com/bupd/timeotter/actions/workflows/test.yml)
[![Lint](https://img.shields.io/github/actions/workflow/status/bupd/timeotter/lint.yml?style=flat-square&logo=github&label=lint)](https://github.com/bupd/timeotter/actions/workflows/lint.yml)
[![Docs](https://img.shields.io/github/actions/workflow/status/bupd/timeotter/docs.yml?style=flat-square&logo=github&label=docs)](https://github.com/bupd/timeotter/actions/workflows/docs.yml)
[![Netlify Status](https://api.netlify.com/api/v1/badges/961e7a52-306c-400d-994f-bdc1ff51e20d/deploy-status)](https://app.netlify.com/projects/timeotter/deploys)

[![GitHub Stars](https://img.shields.io/github/stars/bupd/timeotter?style=flat-square&logo=github)](https://github.com/bupd/timeotter/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/bupd/timeotter?style=flat-square&logo=github)](https://github.com/bupd/timeotter/network/members)
[![GitHub Issues](https://img.shields.io/github/issues/bupd/timeotter?style=flat-square&logo=github)](https://github.com/bupd/timeotter/issues)
[![GitHub PRs](https://img.shields.io/github/issues-pr/bupd/timeotter?style=flat-square&logo=github)](https://github.com/bupd/timeotter/pulls)

[![Homebrew](https://img.shields.io/badge/homebrew-available-orange?style=flat-square&logo=homebrew)](https://github.com/bupd/timeotter#installation)
[![Go Install](https://img.shields.io/badge/go%20install-available-00ADD8?style=flat-square&logo=go)](https://github.com/bupd/timeotter#installation)
[![Platform Linux](https://img.shields.io/badge/platform-linux-FCC624?style=flat-square&logo=linux&logoColor=black)](https://github.com/bupd/timeotter/releases)
[![Platform macOS](https://img.shields.io/badge/platform-macOS-000000?style=flat-square&logo=apple)](https://github.com/bupd/timeotter/releases)

[![Documentation](https://img.shields.io/badge/docs-timeotter.bupd.xyz-blue?style=flat-square&logo=readme)](https://timeotter.bupd.xyz)
[![Made with Go](https://img.shields.io/badge/made%20with-Go-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Google Calendar](https://img.shields.io/badge/Google%20Calendar-API-4285F4?style=flat-square&logo=google-calendar)](https://developers.google.com/calendar)

---

<img align="right" src="artwork/image1.png" width="400" height="410">

A utility that integrates with Google Calendar and helps you run scheduled calendar alarms and execute commands based on them.
In other words `execute commands based on your calendars` 

(ie. Calendar Driven Task Execution)

<!-- <img align="right" src="https://github.com/user-attachments/assets/a7a13a30-be33-445f-b624-7fc93f4a3d00" width="400" height="410"> -->

## Setup Instructions 🪛

Follow the instructions below to set up and use the **Time Otter**:

### Step 1: Generate Google Calendar OAuth Token 📅

Before you can use TimeOtter, you need to authenticate your Google Calendar access.

1. Visit the [Google Developer Console](https://console.developers.google.com/).
2. Create a new project.
3. Enable the Google Calendar API for the project.
4. Create OAuth 2.0 credentials for your project.
5. Download the credentials JSON file.

Use this JSON file to generate the OAuth token for your Google Calendar. You can find instructions on how to generate the OAuth token in [Google's API documentation](https://developers.google.com/calendar/auth).

Once you have the OAuth token, save it to the following location:

```
~/.cal-token.json
```

## Step 2: Configure the Variables in `config.toml` ⚙️

After obtaining the OAuth token, configure the required variables in the `config.toml` file.

Conviniently located at `~/.config/timeotter/config.toml`

Here's an example of what you need to set:

```toml
# config.toml

# Required settings
CalendarID = "your-email@gmail.com"  # Replace with your Google Calendar email or "primary" for default
CmdToExec  = "mpv ~/video.mp4"       # Command to execute when alarm triggers
TokenFile  = "~/.cal-token.json"     # Path to OAuth token (~ is expanded to home directory)

# Optional settings (with defaults)
MaxRes               = 5                                    # Number of events to fetch (min: 1, max: 100, default: 5)
CredentialsFile      = "~/.cal-credentials.json"            # OAuth credentials file path
BackupFile           = "~/.crontab_backup.txt"              # Crontab backup location
TriggerBeforeMinutes = 5                                    # Minutes before event to trigger alarm (default: 5)
CronMarker           = "# custom crons below this can be deleted."  # Delimiter for managed crons
ShowDeleted          = false                                # Include deleted events (default: false)
```

## Step 3: Modify Crontab to Integrate with TimeOtter ⏳

In order for Time Otter to manage your calendar alarms, you need to add the following comment to the **end** of your crontab:

```sh
# custom crons below this can be deleted.
```

This comment marks the entry point for the app to schedule cron jobs. You can customize this marker via the `CronMarker` config option.
### **Do not add any crons below this comment**, as these will be deleted when the app runs.

## 🚨🚨 Important Notes 🚨🚨:

- Before making any changes to your crontab, **take a backup** of your existing cron jobs. You can do this by running:

    ```bash
    crontab -l > crontab-backup.txt
    ```

- After adding the comment, you can proceed with running the application. Time Otter will automatically schedule your calendar-based alarms.

## Step 4: Running the Application 🏄‍♀️

Once you have completed the configuration, you're good to run the application. TimeOtter will fetch events from your Google Calendar and run the corresponding commands when the events are triggered.

Simply execute the program to start syncing your calendar alarms and running the commands you've configured.

```bash
go run time_otter.go
```

### 🧑‍🎤 Running the Application as a Cron Job
Once you have completed the configuration, you're ready to run **Time Otter** as a cron job. This allows **Time Otter** to automatically check your Google Calendar and execute the corresponding commands on a regular basis.

## 👨‍💻 Installation

To install **Time Otter** globally on your system, use the following command:

```bash
go install github.com/bupd/timeotter/cmd/timeotter@latest
```

Make sure you have Go set up correctly in your environment before running the above command.

## 🧠 Step 4.2: Add the Cron Job

Once **Time Otter** is installed, you need to add a cron job that runs **Time Otter** at regular intervals. This ensures that your calendar events are checked, and the configured commands are executed as scheduled.

You can add the following cron job to your crontab:

For running the job every hour:
```sh
0 * * * * timeotter
```

For running the job every 30 minutes:
```sh
*/30 * * * * timeotter
```

> **Important:** Make sure to add this cron job **above** the comment `# custom crons below this can be deleted.` in your crontab.

This will ensure that **Time Otter** will execute as scheduled, while also keeping your custom cron jobs intact.

Once the cron job is set up, **Time Otter** will automatically run at the specified intervals, sync with your Google Calendar, and trigger the corresponding alarms and commands.

You're now all set! Enjoy automated calendar management with **Time Otter**!

> /s Happy abusing Technology 🤩
