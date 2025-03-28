# Yandex Metrika Telegram Bot aka Zakhar

## Overview
This project is a Telegram bot that retrieves and displays statistics from Yandex Metrika. The bot allows authorized users to request daily metrics such as page views, visits, and unique users.

## Features
- Fetches real-time statistics from Yandex Metrika.
- Logs bot interactions using structured JSON logging with `log/slog`.
- Supports access control to limit requests to specific users or channels.
- Uses environment variables for configuration.

## Installation

### Prerequisites
- Go 1.20+
- A Yandex Metrika account and API token
- A Telegram bot token

### Clone the Repository
```sh
git clone https://github.com/rallvesh/zakhar.git
cd zakhar
```

### Configuration
Create a `.env` file in the root directory and set the following variables:
```ini
YANDEX_METRIKA_TOKEN=your_metrika_token
YANDEX_METRIKA_COUNTER_ID=your_counter_id
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
ALLOWED_USER_IDS=your_chat_id      # Comma-separated list of allowed users
ALLOWED_CHAT_ID=your_user_id_list # Allowed chat ID for requests
```

### Build and Run
```sh
go build -o zakhar ./cmd/zakhar/
./zakhar
```

## Usage
- `/start` - Greets the user.
- `/stats` - Fetches and displays today's statistics from Yandex Metrika.

## Logging
The bot uses `slog` for structured logging in JSON format. All commands are logged, but regular messages are ignored.

## License
MIT License. See `LICENSE` for details.
