# TimeOtter

Google Calendar integration utility that syncs calendar events to cron jobs.

## Style Preferences

- Use `##` for headings, `###` for subheadings
- No `**` bold or long dashes
- File names in backticks like `config.go`
- Use GitHub file links when referencing code locations
- Prefer conciseness over grammar
- No tabs in markdown

## Documentation Guidelines

- Update README.md when logic or config options change
- Keep docs in sync with code changes

## Project Structure

```
cmd/timeotter/main.go    - Entry point
pkg/config/config.go     - Viper-based TOML config
pkg/oauth/oauth.go       - Google OAuth2 flow
pkg/calendar/calendar.go - Event parsing, cron conversion
pkg/cron/cron.go         - Crontab manipulation
pkg/utils/utils.go       - Shell command execution
```

## Config Location

`~/.config/timeotter/config.toml`

## Key Files

- Credentials: configurable via `CredentialsFile` (default: `~/.cal-credentials.json`)
- Token: configurable via `TokenFile`
- Cron backup: configurable via `BackupFile` (default: `~/.crontab_backup.txt`)

## Config Options

All options in `pkg/config/config.go`:

- `CalendarID` - Google Calendar ID
- `CmdToExec` - Command to run on event trigger
- `MaxRes` - Max events to fetch (1-100, default: 5)
- `TokenFile` - OAuth token path
- `CredentialsFile` - OAuth credentials path
- `BackupFile` - Crontab backup location
- `TriggerBeforeMinutes` - Minutes before event to trigger alarm (default: 5)
- `CronMarker` - Delimiter comment for managed crons
- `ShowDeleted` - Include deleted events (default: false)

## Build

```bash
go build -o timeotter cmd/timeotter/main.go
```

## Dependencies

- spf13/viper (config)
- google.golang.org/api (calendar)
- golang.org/x/oauth2 (auth)
