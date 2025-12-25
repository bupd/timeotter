---
title: Getting Started
description: Get up and running with TimeOtter in minutes
---

TimeOtter integrates with Google Calendar to execute commands based on your calendar events. Follow these steps to get started.

## Prerequisites

- Go 1.21+ (for installation via `go install`)
- A Google account with Calendar access
- Linux/macOS with cron support

## Quick Setup

### 1. Install TimeOtter

```bash
go install github.com/bupd/timeotter/cmd/timeotter@latest
```

Or via Homebrew (macOS):
```bash
brew install bupd/tap/timeotter
```

### 2. Set up Google OAuth

Follow the [OAuth Setup guide](/oauth-setup) to create credentials and generate your token.

### 3. Configure TimeOtter

Create `~/.config/timeotter/config.toml`:

```toml
CalendarID = "your-email@gmail.com"
CmdToExec  = "mpv ~/alarm.mp3"
TokenFile  = "~/.cal-token.json"
```

### 4. Add to Cron

Add TimeOtter to your crontab to run periodically:

```bash
crontab -e
```

Add this line (runs every 30 minutes):
```
*/30 * * * * timeotter
```

See [Cron Setup](/cron-setup) for more details.

## Next Steps

- [Installation options](/installation) - Different ways to install
- [Configuration reference](/configuration) - All config options
- [Troubleshooting](/troubleshooting) - Common issues and fixes
