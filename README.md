# timeotter


<!-- <img align="right" src="https://github.com/user-attachments/assets/a7a13a30-be33-445f-b624-7fc93f4a3d00" width="400" height="410"> -->

<img align="right" src="https://github.com/user-attachments/assets/073bf189-3c47-4ec3-bc60-97ed9d288b2a" width="400" height="410">

## Setup Instructions

Follow the instructions below to set up and use the **Time Otter**:

### Step 1: Generate Google Calendar OAuth Token

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

### Step 2: Configure the Variables in `config.toml`

After obtaining the OAuth token, configure the required variables in the `config.toml` file.

Conviniently located at `~/.config/timeotter/config.toml`

Here's an example of what you need to set:

```toml
# config.toml

CalendarID = "your-email@gmail.com"  # Replace with your Google Calendar email address or primary for default calendar
CmdToExec  = "mpv ~/video.mp4"       # Replace with the command you wish to execute when the alarm triggers
MaxRes     = 2                        # Number of results to fetch from Google Calendar API (adjust as necessary) advised to set below 10 for safety
TokenFile  = "~/.cal-token.json"      # Path to your OAuth token
```

### Step 3: Modify Crontab to Integrate with TimeOtter

In order for Time Otter to manage your calendar alarms, you need to add the following comment to the ==**end**== your crontab:

```sh
# custom crons below this can be deleted.
```

This comment marks the entry point for the app to schedule cron jobs.
### **Do not add any crons below this comment**, as these will be deleted when the app runs.

#### Important Notes:

- Before making any changes to your crontab, **take a backup** of your existing cron jobs. You can do this by running:

    ```bash
    crontab -l > crontab-backup.txt
    ```

- After adding the comment, you can proceed with running the application. Time Otter will automatically schedule your calendar-based alarms.

### Step 4: Running the Application

Once you have completed the configuration, you're good to run the application. TimeOtter will fetch events from your Google Calendar and run the corresponding commands when the events are triggered.

Simply execute the program to start syncing your calendar alarms and running the commands you've configured.

```bash
go run time_otter.go
```

### Running the Application as a Cron Job

Once you have completed the configuration, you're ready to run **Time Otter** as a cron job. This allows **Time Otter** to automatically check your Google Calendar and execute the corresponding commands on a regular basis.

#### Step 4.1: Install Time Otter

To install **Time Otter** globally on your system, use the following command:

```bash
go install github.com/bupd/timeotter@latest
```

Make sure you have Go set up correctly in your environment before running the above command.

#### Step 4.2: Add the Cron Job

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

> Happy abusing Technology
