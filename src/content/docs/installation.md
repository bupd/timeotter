---
title: Installation
description: Different ways to install TimeOtter on your system
---

## Go Install (Recommended)

If you have Go installed, this is the quickest way:

```bash
go install github.com/bupd/timeotter/cmd/timeotter@latest
```

Make sure `$GOPATH/bin` is in your `PATH`.

## Homebrew (macOS)

```bash
brew install bupd/tap/timeotter
```

## Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/bupd/timeotter/releases):

1. Download the appropriate archive for your platform
2. Extract the binary
3. Move it to a directory in your `PATH`

```bash
# Example for Linux amd64
tar -xzf timeotter_Linux_x86_64.tar.gz
sudo mv timeotter /usr/local/bin/
```

## Build from Source

```bash
git clone https://github.com/bupd/timeotter.git
cd timeotter
go build -o timeotter ./cmd/timeotter
sudo mv timeotter /usr/local/bin/
```

## Verify Installation

```bash
timeotter --help
```

## Next Steps

After installation, proceed to [OAuth Setup](/oauth-setup) to configure Google Calendar access.
