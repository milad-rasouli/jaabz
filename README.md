# Jaabz Job Scraper

Jaabz Job Scraper is a Go application that periodically fetches job listings from [Jaabz](https://jaabz.com), checks for duplicates using Redis, and posts non-duplicate jobs to a Telegram channel. The application runs in a Dockerized environment with Redis for persistence and uses the Telegram Bot API for notifications.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
    - [Clone the Repository](#clone-the-repository)
    - [Configure Environment Variables](#configure-environment-variables)
    - [Build and Run with Docker](#build-and-run-with-docker)
    - [Run Locally (Optional)](#run-locally-optional)
- [Usage](#usage)
- [Directory Structure](#directory-structure)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Features
- Fetches job listings from Jaabz every 60 seconds.
- Uses Redis to track and skip duplicate jobs based on their visit links.
- Posts non-duplicate jobs to a Telegram channel with formatted details (title, company, work status, location, skills, apply link).
- Batches multiple jobs into a single Telegram message to reduce API calls and respect rate limits.
- Handles Telegram rate limits with retries and delays.
- Runs in a Dockerized environment for easy deployment.
- Comprehensive logging for debugging and monitoring.

## Prerequisites
- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
- A Telegram bot token from [@BotFather](https://t.me/BotFather).
- A Telegram channel (public or private) where the bot is an admin with **Post Messages** permission.
- Redis (provided via Docker) for duplicate checking.
- Go 1.22 or later (if running locally without Docker).

## Setup

### Clone the Repository
```bash
git clone https://github.com/milad-rasouli/jaabz.git
cd jaabz
```

### Configure Environment Variables
1. Copy the example environment file:
   ```bash
   cp env.example .env
   ```
2. Edit `.env` with your configuration:
   ```env
   APP_NAME=jaabz
   ENVIRONMENT=development
   JAABZ_HOST=https://jaabz.com/jobs?job_state=&keyword=golang&category=&skill=&country=
   REDIS_HOST=jaabz-redis:6379
   TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
   TELEGRAM_CHANNEL_ID=@JaabzJobs
   ```
    - **APP_NAME**: Application name (default: `jaabz`).
    - **ENVIRONMENT**: Environment (e.g., `development`, `production`).
    - **JAABZ_HOST**: URL for fetching Jaabz jobs (update query parameters as needed).
    - **REDIS_HOST**: Redis host and port (use `jaabz-redis:6379` for Docker).
    - **TELEGRAM_BOT_TOKEN**: Obtain from [@BotFather](#obtaining-a-telegram-bot-token).
    - **TELEGRAM_CHANNEL_ID**: Channel handle (e.g., `@JaabzJobs`) or private ID (e.g., `-1001234567890`).

3. **Obtaining a Telegram Bot Token**:
    - Message `@BotFather` on Telegram.
    - Send `/newbot`, set a name (e.g., `Jaabz Job Bot`) and username (e.g., `@JaabzJobBot`).
    - Copy the provided token (e.g., `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`).

4. **Setting Up the Telegram Channel**:
    - Create a public channel (e.g., `@JaabzJobs`) or private channel.
    - For private channels, get the ID by forwarding a message to `@GetIDsBot`.
    - Add the bot as an admin:
        - Channel settings > **Administrators** > **Add Administrator** > `@JaabzJobBot` > Grant **Post Messages** permission.

### Build and Run with Docker
1. Build and start the services:
   ```bash
   docker-compose up -d --build
   ```
    - This starts the `jaabz-app` (Go application) and `jaabz-redis` (Redis) containers.
    - The application is built from the `Dockerfile` and runs the compiled binary at `/app/build`.

2. View logs:
   ```bash
   docker logs jaabz-app
   ```
    - Look for:
      ```
      time=2025-05-15T15:03:06.602Z level=INFO msg="Telegram bot initialized" package=telegram bot_username=@JaabzJobBot
      time=2025-05-15T15:03:06.603Z level=INFO msg="Fetched jobs" service=jaabz method=processJobs count=3
      time=2025-05-15T15:03:06.605Z level=DEBUG msg="Posted batch of jobs to Telegram" service=jaabz method=processJobs job_count=2
      ```

3. Stop the services:
   ```bash
   docker-compose down
   ```

### Run Locally (Optional)
1. Install Go dependencies:
   ```bash
   go mod download
   ```
2. Update `.env` to use local Redis:
   ```env
   REDIS_HOST=localhost:6379
   ```
3. Run Redis locally (if not using Docker):
   ```bash
   docker run -d --name redis -p 6379:6379 redis:8.0-M02-alpine
   ```
4. Build and run:
   ```bash
   make run
   ```
    - The `Makefile` builds the binary to `bin/jaabz` and runs it.

## Usage
- The application fetches jobs from `JAABZ_HOST` every 60 seconds.
- It checks for duplicates using Redis (based on job visit links).
- Non-duplicate jobs are batched into a single Telegram message (up to 4096 characters) and posted to the specified channel.
- Example Telegram post:
  ```
  *New Job Postings*

  *Title*: MongoDB Database Administrator
  *Company*: Tech Corp
  *Work Status*: Remote
  *Location*: Remote
  *Skills*: MongoDB, Linux
  *Apply*: [Link](https://jaabz.com/jobs/81287-mongodb-database-administrator)

  *Title*: Software Engineer
  *Company*: InnoSoft
  *Work Status*: Full-time
  *Location*: San Francisco
  *Skills*: Go, Docker
  *Apply*: [Link](https://jaabz.com/jobs/81217-software-engineer)
  ```
- Rate limits are handled with retries (up to 3 attempts) and delays based on Telegram’s `retry_after` response.

## Directory Structure
```
jaabz/
├── bin/                    # Compiled binaries
├── cmd/                    # Main application entry point
│   └── main.go
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker build instructions
├── env.example             # Example environment file
├── go.mod                  # Go module dependencies
├── go.sum                  # Go dependency checksums
├── internal/               # Internal packages
│   ├── infra/              # Infrastructure (e.g., godotenv, redis)
│   ├── repo/               # Repositories (e.g., telegram, jaabz, duplicate)
│   └── service/            # Business logic
├── Makefile                # Build and run commands
└── README.md               # Project documentation
```

## Environment Variables
| Variable             | Description                                      | Example                                      |
|----------------------|--------------------------------------------------|----------------------------------------------|
| `APP_NAME`           | Application name                                 | `jaabz`                                      |
| `ENVIRONMENT`        | Environment (e.g., development, production)       | `development`                                |
| `JAABZ_HOST`         | URL for fetching Jaabz jobs                      | `https://jaabz.com/jobs?...`                 |
| `REDIS_HOST`         | Redis host and port                              | `jaabz-redis:6379` (Docker) or `localhost:6379` |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token from @BotFather               | `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11` |
| `TELEGRAM_CHANNEL_ID`| Telegram channel ID (public or private)           | `@JaabzJobs` or `-1001234567890`             |

## Troubleshooting
1. **Telegram "Too Many Requests"**:
    - The application retries up to 3 times with delays. If persistent, reduce the number of jobs fetched or increase the fetch interval in `service.go`.
    - Check logs:
      ```
      level=WARN msg="Rate limit hit, retrying" package=telegram attempt=1 retry_after_seconds=38
      ```

2. **Telegram "Not Found" or Invalid Token**:
    - Verify `TELEGRAM_BOT_TOKEN` in `.env`.
    - Test:
      ```bash
      curl "https://api.telegram.org/bot<your-token>/getMe"
      ```
    - Regenerate token via `@BotFather` if needed.

3. **Invalid Channel ID**:
    - Ensure `TELEGRAM_CHANNEL_ID` is `@ChannelName` or `-100NNNNNNNNNN`.
    - For private channels, use `@GetIDsBot`.

4. **Redis Connection Issues**:
    - Confirm Redis is running:
      ```bash
      docker ps  # Check jaabz-redis
      redis-cli -h localhost -p 6378 ping  # If running locally
      ```
    - Update `REDIS_HOST` if using a non-standard port (e.g., `localhost:6378`).

5. **Logs**:
    - Check `docker logs jaabz-app` or `bin/jaabz` output for errors.
    - Enable debug logging in `main.go` (`slog.LevelDebug`).

## Contributing
1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/xyz`).
3. Commit changes (`git commit -m "Add feature xyz"`).
4. Push to the branch (`git push origin feature/xyz`).
5. Open a pull request.

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.


### Notes on the README
- **Project Overview**: Describes the application’s purpose (scraping Jaabz jobs, deduplicating with Redis, posting to Telegram).
- **Docker Integration**: Details the Docker Compose setup, including Redis and the Go app, with port mapping (`6378:6379`) from your configuration.
- **Environment Variables**: Lists all variables from `docker-compose.yml` and explains how to obtain `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHANNEL_ID`.
- **Setup Instructions**: Covers both Docker and local execution, including `make run` and Redis setup.
- **Usage**: Explains the application’s behavior and shows an example Telegram post.
- **Directory Structure**: Reflects the provided `ls` output (`bin/`, `cmd/`, etc.).
- **Troubleshooting**: Addresses common issues (rate limits, invalid tokens, Redis), referencing your recent errors.
- **Assumptions**:
    - The `Dockerfile` builds the Go binary to `/app/build` (per `entrypoint` in `docker-compose.yml`).
    - The `JAABZ_HOST` and Telegram variables are left blank in the example `.env` to be filled by the user.
    - Redis uses port `6378` externally (mapped to `6379` internally) as per your `docker-compose.yml`.

### Additional Setup Notes
1. **Dockerfile**: Ensure your `Dockerfile` compiles the Go binary to `/app/build`. Example:
   ```dockerfile
   FROM golang:1.22-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o build ./cmd/.

   FROM alpine:latest
   WORKDIR /app
   COPY --from=builder /app/build .
   CMD ["/app/build"]
   ```

2. **env.example**: Create an `env.example` file to match your `docker-compose.yml`:
   ```env
   APP_NAME=jaabz
   ENVIRONMENT=development
   JAABZ_HOST=
   REDIS_HOST=jaabz-redis:6379
   TELEGRAM_CHANNEL_ID=
   TELEGRAM_BOT_TOKEN=
   ```