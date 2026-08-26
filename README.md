# Notification System

A scalable notification service built with Go, the Chi router, PostgreSQL, and an asynchronous background worker with retry handling. The project exposes a REST API for creating and querying notifications while processing delivery jobs in the background to improve throughput and resilience.

## Overview

This service is designed to support notification creation and processing at scale using a clean layered architecture. It allows clients to create records for user notifications, query them with filters and pagination, and asynchronously process pending notifications for delivery. The system persists notification state in PostgreSQL and uses a background worker to transition records from `Pending` to `Processing` and finally to either `Sent` or `Failed` depending on delivery outcome.

## Key Features

- RESTful API for notification management:
  - Create a notification
  - List notifications by user with pagination and unread filtering
  - Retrieve a single notification by ID and mark it as read
- Asynchronous background worker for processing pending notifications
- Resilient retry mechanism with status transitions:
  - `Pending` -> `Processing` -> `Sent`
  - `Pending` -> `Processing` -> `Failed`
- Clean architecture separation:
  - HTTP handlers
  - Repository layer
  - Provider abstraction
  - Background worker orchestration
- Dockerized local development setup with PostgreSQL

## Tech Stack

- Go
- Chi Router (`github.com/go-chi/chi/v5`)
- PostgreSQL
- `pg`/`lib/pq` database driver
- Docker & Docker Compose

## Project Structure

```text
notification/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── main.go
├── application/
│   ├── app.go
│   └── routes.go
├── database/
│   └── db.go
├── handler/
│   └── notification.go
├── model/
│   └── notification.go
├── pkg/
│   └── provider/
│       ├── email.go
│       └── email_test.go
├── repository/
│   ├── notification.go
│   └── notification_worker_repo.go
├── worker/
│   └── notification_worker.go
└── README.md
```

## Environment Variables

The application reads its configuration from environment variables, and it falls back to sensible defaults when values are not provided.

| Variable      | Description              | Example             |
| ------------- | ------------------------ | ------------------- |
| `PORT`        | HTTP server port         | `4000`              |
| `DB_HOST`     | PostgreSQL host          | `localhost` or `db` |
| `DB_PORT`     | PostgreSQL port          | `5432`              |
| `DB_USER`     | PostgreSQL username      | `postgres`          |
| `DB_PASSWORD` | PostgreSQL password      | `postgrespassword`  |
| `DB_NAME`     | PostgreSQL database name | `notification_db`   |
| `SSL_MODE`    | PostgreSQL SSL mode      | `disable`           |

> The app initializes the database connection using `godotenv` and defaults `SSL_MODE` to `disable` if it is not set.

## Getting Started

### Prerequisites

- Go 1.27 or later
- Docker and Docker Compose
- PostgreSQL (if running outside Docker)

### Run with Docker Compose

From the project root, start the app and PostgreSQL with:

```bash
docker-compose up --build
```

This will:

- start a PostgreSQL container
- build the Go application container
- expose the API on port `4000`
- initialize the `notifications` table automatically if it does not exist

### Run Without Docker

1. Make sure PostgreSQL is running and accessible.
2. Create a local `.env` file or export environment variables:

```bash
export PORT=4000
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgrespassword
export DB_NAME=notification_db
export SSL_MODE=disable
```

3. Run the application:

```bash
go run .
```

## API Endpoints Reference

### Health Check

| Method | Route | Description                                            |
| ------ | ----- | ------------------------------------------------------ |
| `GET`  | `/`   | Returns HTTP 200 to confirm the application is running |

### Notification Management

| Method | Route                                         | Description                                  |
| ------ | --------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------------ |
| `POST` | `/notification`                               | Create a new notification                    |
| `GET`  | `/notification?user_id={id}&unread_only={true | false}`                                      | List notifications for a user with pagination and unread filtering |
| `GET`  | `/notification/{id}`                          | Get a notification by ID and mark it as read |

### Request and Query Details

#### Create Notification

`POST /notification`

Request body:

```json
{
  "user_id": "user-123",
  "type": "EMAIL",
  "message": "Welcome to our platform!"
}
```

Allowed notification types:

- `EMAIL`
- `SMS`
- `PUSH`

#### List Notifications

`GET /notification?user_id={id}&unread_only={true|false}`

Supported query parameters:

- `user_id` (required)
- `page` (optional, defaults to `1`)
- `limit` (optional, defaults to `10`, max `100`)
- `unread_only` (optional, accepts `true` or `false`)

Example:

```bash
curl "http://localhost:4000/notification?user_id=user-123&unread_only=true&page=1&limit=10"
```

#### Get Single Notification

`GET /notification/{id}`

This endpoint fetches the notification by ID and sets `is_read` to `true` in the database.

## Worker & Retry Architecture

The application starts a notification worker when the service boots. The worker uses a ticker and polls the database every few seconds to fetch pending notifications.

Processing flow:

1. Query notifications where the status is `Pending` or `Failed` and `retry_count < 3`.
2. Mark each matching record as `Processing`.
3. Deliver the notification through the `EmailProvider` abstraction.
4. If delivery succeeds, update the status to `Sent`.
5. If delivery fails, increment `retry_count` and keep the notification in `Pending` until it reaches the configured max retries.
6. When the retry limit is reached, set the record to `Failed`.

The worker is implemented in the `worker` package, while the provider abstraction is defined under `pkg/provider`. The current implementation uses a mock email provider, which simulates email delivery and can randomly fail to mimic real provider behavior.

## Notes

- The app automatically creates the `notifications` table on startup if it does not already exist.
- The background worker is a lightweight polling worker and uses a simple retry model suited for local development and service orchestration.
- The service is intentionally structured to separate concerns for easier extension with actual email/SMS/push integrations.


