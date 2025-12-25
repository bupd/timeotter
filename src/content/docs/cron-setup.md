---
title: Cron Setup
description: Run TimeOtter automatically as a cron job
---

TimeOtter works best when run periodically via cron. This guide explains how to set it up.

## How It Works

1. TimeOtter fetches your upcoming calendar events
2. It creates cron entries for each event
3. When an event time arrives, cron executes your configured command

## Setting Up Your Crontab

### Step 1: Add the Marker Comment

First, add the marker comment to your crontab. This tells TimeOtter where to manage entries:

```bash
crontab -e
```

Add this line at the end:

```
# custom crons below this can be deleted.
```

:::caution
Do not add any cron jobs below this line. TimeOtter will delete them when it runs.
:::

### Step 2: Add TimeOtter Cron Job

Add a cron entry to run TimeOtter periodically. Place this **above** the marker comment:

```bash
# Run TimeOtter every 30 minutes
*/30 * * * * timeotter

# custom crons below this can be deleted.
```

### Example Complete Crontab

```bash
# Your existing cron jobs
0 0 * * * /path/to/backup-script.sh
0 9 * * 1 /path/to/weekly-report.sh

# Run TimeOtter every hour
0 * * * * timeotter

# custom crons below this can be deleted.
# (TimeOtter manages everything below this line)
```

## Backup Your Crontab

Before making changes, back up your existing crontab:

```bash
crontab -l > ~/crontab-backup.txt
```

TimeOtter also creates automatic backups at the path specified by `BackupFile` in your config.

## Recommended Schedules

| Frequency | Cron Expression | Use Case |
|-----------|-----------------|----------|
| Every 30 min | `*/30 * * * *` | Frequent calendar updates |
| Every hour | `0 * * * *` | Standard usage |
| Every 2 hours | `0 */2 * * *` | Less frequent updates |
| Every 6 hours | `0 */6 * * *` | Minimal resource usage |

## Verifying Setup

Check that your crontab is correctly configured:

```bash
crontab -l
```

You should see:
1. The TimeOtter cron job (above the marker)
2. The marker comment
3. Any calendar-generated entries (below the marker)

## Troubleshooting

### TimeOtter Not Running

Check cron logs:

```bash
grep CRON /var/log/syslog
```

### Command Not Found

Ensure the full path to timeotter is used if not in PATH:

```bash
0 * * * * /home/user/go/bin/timeotter
```

### No Events Appearing

1. Verify your config file is correct
2. Check that you have upcoming calendar events
3. Run `timeotter` manually to see any errors
