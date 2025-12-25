---
title: Troubleshooting
description: Common issues and solutions for TimeOtter
---

## Authentication Issues

### "Token has expired"

The OAuth token needs to be refreshed. Delete the token and re-authenticate:

```bash
rm ~/.cal-token.json
timeotter
```

### "Access Denied" during OAuth

1. Make sure you added your email as a test user in Google Cloud Console
2. Click **Advanced** → **Go to TimeOtter (unsafe)** to proceed

### "Invalid credentials"

Verify your credentials file exists and is valid:

```bash
cat ~/.cal-credentials.json
```

If missing or corrupt, re-download from Google Cloud Console.

## Calendar Issues

### No Events Found

1. Check that `CalendarID` is correct in your config
2. Verify you have upcoming events in your calendar
3. Try using `"primary"` as the CalendarID

### Wrong Calendar

Find your calendar ID:
1. Open Google Calendar
2. Click the three dots next to your calendar
3. Select **Settings and sharing**
4. Scroll to **Integrate calendar**
5. Copy the **Calendar ID**

## Cron Issues

### TimeOtter Not Running Automatically

1. Check cron is running:
   ```bash
   systemctl status cron
   ```

2. Verify the cron entry:
   ```bash
   crontab -l
   ```

3. Check cron logs:
   ```bash
   grep CRON /var/log/syslog
   ```

### Events Not Triggering

1. Verify the marker comment is in your crontab
2. Run `timeotter` manually to see generated entries
3. Check that event times are in the future

### "Command not found" in Cron

Cron uses a limited PATH. Use the full path to timeotter:

```bash
0 * * * * /home/user/go/bin/timeotter
```

Or to your command:

```toml
CmdToExec = "/usr/bin/mpv ~/alarm.mp3"
```

## Configuration Issues

### Config File Not Found

Create the config directory and file:

```bash
mkdir -p ~/.config/timeotter
touch ~/.config/timeotter/config.toml
```

### Invalid TOML Syntax

Check your config for common issues:

- Strings must be quoted: `CalendarID = "email@gmail.com"`
- No trailing commas
- Use `=` not `:`

Validate your TOML:
```bash
cat ~/.config/timeotter/config.toml
```

## Command Issues

### Command Not Executing

1. Test the command manually:
   ```bash
   mpv ~/alarm.mp3
   ```

2. Check file permissions
3. Use absolute paths in `CmdToExec`

### Audio/Video Not Playing

Cron runs without a display. For GUI commands:

```toml
CmdToExec = "DISPLAY=:0 mpv ~/alarm.mp3"
```

For audio:

```toml
CmdToExec = "paplay ~/alarm.wav"
```

## Getting Help

If you're still stuck:

1. Check [GitHub Issues](https://github.com/bupd/timeotter/issues)
2. Open a new issue with:
   - Your config file (redact sensitive data)
   - Error messages
   - Steps to reproduce
