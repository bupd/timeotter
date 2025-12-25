---
title: Configuration
description: Complete configuration reference for TimeOtter
---

TimeOtter is configured via a TOML file located at `~/.config/timeotter/config.toml`.

## Example Configuration

```toml
# Required settings
CalendarID = "your-email@gmail.com"
CmdToExec  = "mpv ~/alarm.mp3"
TokenFile  = "~/.cal-token.json"

# Optional settings
MaxRes               = 5
CredentialsFile      = "~/.cal-credentials.json"
BackupFile           = "~/.crontab_backup.txt"
TriggerBeforeMinutes = 5
CronMarker           = "# custom crons below this can be deleted."
ShowDeleted          = false
```

## Required Settings

### CalendarID

The Google Calendar ID to fetch events from.

```toml
CalendarID = "your-email@gmail.com"
```

- Use `"primary"` for your default calendar
- Or use your email address
- For shared calendars, use the calendar ID from Calendar settings

### CmdToExec

The shell command to execute when an event triggers.

```toml
CmdToExec = "mpv ~/alarm.mp3"
```

Examples:
- `"notify-send 'Meeting starting!'"` - Desktop notification
- `"mpv ~/sounds/alarm.mp3"` - Play audio
- `"/path/to/script.sh"` - Run a script

### TokenFile

Path to the OAuth token file.

```toml
TokenFile = "~/.cal-token.json"
```

The `~` is expanded to your home directory.

## Optional Settings

### MaxRes

Number of upcoming events to fetch.

```toml
MaxRes = 5
```

- **Default:** 5
- **Range:** 1-100

### CredentialsFile

Path to OAuth client credentials.

```toml
CredentialsFile = "~/.cal-credentials.json"
```

- **Default:** `~/.cal-credentials.json`

### BackupFile

Where to store crontab backups.

```toml
BackupFile = "~/.crontab_backup.txt"
```

- **Default:** `~/.crontab_backup.txt`

### TriggerBeforeMinutes

How many minutes before an event to trigger the command.

```toml
TriggerBeforeMinutes = 5
```

- **Default:** 5
- Set to `0` to trigger at event start time

### CronMarker

Comment that marks where TimeOtter manages cron entries.

```toml
CronMarker = "# custom crons below this can be deleted."
```

- **Default:** `# custom crons below this can be deleted.`
- Place this in your crontab to mark the start of TimeOtter-managed entries
- Everything below this line will be replaced by TimeOtter

### ShowDeleted

Whether to include deleted events.

```toml
ShowDeleted = false
```

- **Default:** `false`

## Environment Variables

TimeOtter also respects the following environment variables:

| Variable | Description |
|----------|-------------|
| `HOME` | Used for `~` expansion in paths |

## Config File Location

TimeOtter looks for configuration in:

1. `~/.config/timeotter/config.toml`

Create the directory if it doesn't exist:

```bash
mkdir -p ~/.config/timeotter
```
