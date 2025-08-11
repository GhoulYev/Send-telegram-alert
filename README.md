# Telegram Alert Sender for Linux

#### A lightweight CLI tool for sending Telegram alerts from Linux systems. Perfect for cron jobs, system monitoring, and custom notifications.

## Features
* 📨 Send Telegram messages with a simple CLI command

* 🔑 Secure credential storage

* ⏱️ Perfect for cron jobs and automated scripts

* 📦 Single binary with no dependencies

## Installation

### Prerequisites

* Go 1.16+ (for building from source)

* Telegram account (for bot setup)


## Build from Source

```bash
# Clone repository
git clone https://github.com/GhoulYev/Send-telegram-alert.git
cd sta

# Build binary
go build -o sta

# Install system-wide
sudo mv sta /usr/local/bin/
```

## Configuration

### First-Time Setup

#### 1. Run the tool to generate a configuration template:
```bash
sta -m "Test message"
```

#### 2. Fill in your credentials in the generated config file:
```bash
nano ~/.config/sta/config.json
```

#### Example configuration:
```json
{
  "token": "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghi",
  "chat_id": "123456789",
  "prefix": "[My Server]"
}
```

#### 4. Get your Telegram credentials:
* Create a bot with [@BotFather](https://t.me/BotFather) to get your token

* Get your chat_id with [@userinfobot](https://t.me/userinfobot)

## Custom Configuration Path

#### Use the -c flag to specify a custom config location:
```bash
sta -c /path/to/config.json -m "Custom config test"
```

## Usage

### Basic Notification

```bash
sta -m "Server started successfully"
```

### System Alerts

```bash
# Disk space alert
sta -m "Disk space critical: $(df -h / | awk 'NR==2{print $5}') used"
```

### With Variables

```bash
# Send current load average
sta -m "Current load: $(uptime | awk -F'[a-z]:' '{print $2}')"
```

### In Scripts

```bash
#!/bin/bash

backup_result=$(rsync -a /source /backup)
if [ $? -eq 0 ]; then
    sta -m "Backup completed successfully"
else
    sta -m "Backup failed: $backup_result"
fi
```

## Cron Integration

### Example crontab entries:

```bash
# Daily system status report
0 9 * * * /usr/local/bin/sta -m "Daily report: $(uptime)"

# Disk space monitoring (every 6 hours)
0 */6 * * * [ $(df -h / | awk 'NR==2 {print $5}' | tr -d '%') -gt 90 ] && /usr/local/bin/sta -m "Disk space critical!"

# Weekly updates
0 0 * * 1 /usr/local/bin/sta -m "Starting weekly updates" && pacman -Syu
```

## Troubleshooting
### Common Issues:

* Error loading configuration: Check config file permissions and syntax

* Token or Chat ID missing: Verify your config contains valid credentials

* API error: 400 Bad Request: Usually indicates an invalid chat ID

* API error: 401 Unauthorized: Verify your bot token

