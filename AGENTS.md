# TimeOtter

Google Calendar integration utility that syncs calendar events to cron jobs.

## Style Preferences

- Use `##` for headings, `###` for subheadings
- No `**` bold or long dashes
- File names in backticks like `config.go`
- Use GitHub file links when referencing code locations
- Prefer conciseness over grammar
- No tabs in markdown

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

- Credentials: `~/.cal-credentials.json`
- Token: configurable via `TokenFile`
- Cron backup: `~/.crontab_backup.txt`

## Hardcoded Values to Make Configurable

- 5-minute event advance warning in `calendar/calendar.go:54`
- Credentials file path in `main.go:34`
- Cron delimiter string in `cron/cron.go:62`
- Backup file location in `cron/cron.go:58`
- API params (ShowDeleted, SingleEvents, OrderBy) in `main.go:74-75`

## Build

```bash
go build -o timeotter cmd/timeotter/main.go
```

## Dependencies

- spf13/viper (config)
- google.golang.org/api (calendar)
- golang.org/x/oauth2 (auth)
