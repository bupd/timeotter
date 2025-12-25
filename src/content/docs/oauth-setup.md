---
title: OAuth Setup
description: Set up Google Calendar API credentials for TimeOtter
---

TimeOtter needs OAuth credentials to access your Google Calendar. Follow these steps to set up authentication.

## Step 1: Create a Google Cloud Project

1. Visit the [Google Cloud Console](https://console.cloud.google.com/)
2. Click **Select a project** → **New Project**
3. Name it (e.g., "TimeOtter") and click **Create**

## Step 2: Enable the Calendar API

1. Go to **APIs & Services** → **Library**
2. Search for "Google Calendar API"
3. Click on it and press **Enable**

## Step 3: Configure OAuth Consent Screen

1. Go to **APIs & Services** → **OAuth consent screen**
2. Select **External** and click **Create**
3. Fill in the required fields:
   - App name: TimeOtter
   - User support email: Your email
   - Developer contact: Your email
4. Click **Save and Continue**
5. On Scopes, click **Add or Remove Scopes**
6. Add: `https://www.googleapis.com/auth/calendar.readonly`
7. Click **Save and Continue**
8. Add your email as a test user
9. Click **Save and Continue**

## Step 4: Create OAuth Credentials

1. Go to **APIs & Services** → **Credentials**
2. Click **Create Credentials** → **OAuth client ID**
3. Application type: **Desktop app**
4. Name: TimeOtter
5. Click **Create**
6. Download the JSON file

## Step 5: Save Credentials

Save the downloaded JSON file:

```bash
mv ~/Downloads/client_secret_*.json ~/.cal-credentials.json
```

## Step 6: Generate Token

Run TimeOtter once to generate the OAuth token:

```bash
timeotter
```

This will:
1. Open your browser for authentication
2. Ask you to authorize TimeOtter
3. Save the token to `~/.cal-token.json`

## File Locations

| File | Default Path | Description |
|------|--------------|-------------|
| Credentials | `~/.cal-credentials.json` | OAuth client credentials |
| Token | `~/.cal-token.json` | Access/refresh token |

Both paths can be customized in your [configuration](/configuration).

## Troubleshooting

### "Access Denied" or "App Not Verified"

This is normal for development. Click **Advanced** → **Go to TimeOtter (unsafe)** to proceed.

### Token Expired

Delete the token file and run TimeOtter again to re-authenticate:

```bash
rm ~/.cal-token.json
timeotter
```
