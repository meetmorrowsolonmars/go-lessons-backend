# Education pet project

[//]: # (Project description coming soon...)

## Prerequisites

- Go 1.24.0 or higher
- Docker and Docker Compose
- PostgreSQL (via Docker)

## Setup and Installation

### Install Tools

Install the Goose migration tool locally in the project's `bin` directory and verify the installation:

```shell
GOBIN="$(pwd)/bin" go install github.com/pressly/goose/v3/cmd/goose@v3.24.3
GOBIN="$(pwd)/bin" go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5
PATH="$(pwd)/bin:${PATH}" goose --version
PATH="$(pwd)/bin:${PATH}" minimock --version

PATH="$(pwd)/bin:${PATH}" go generate ./...

go list ./...

go test ./...

go test -tags integration,local ./...
go test -tags integration,ci ./...
```

```shell
go test ./internal/... ./pkg/...

./scripts/integration-test.sh
```

**Alternative installation method:** use `go get -tool github.com/pressly/goose/v3/cmd/goose@v3.24.3` and
`go install tool` as an alternative way to install tools.

### Database Migration Management

Generate a new SQL migration file for creating database tables or schema changes:

```shell
PATH="$(pwd)/bin:${PATH}" goose -dir=migrations create create_user_table sql
```

### Infrastructure Setup

Launch the PostgreSQL database using Docker Compose in detached mode:

```shell
docker compose up -d
```

Run all pending migrations to update the database schema:

```shell
export GOOSE_DBSTRING="dbname=postgres user=postgres password=postgres host=127.0.0.1 port=25432 sslmode=disable"

PATH="$(pwd)/bin:${PATH}" goose -dir=migrations postgres status
PATH="$(pwd)/bin:${PATH}" goose -dir=migrations postgres up
```

### Start Application

```shell
export DATABASE_URL="dbname=postgres user=postgres password=postgres host=127.0.0.1 port=25432 sslmode=disable"
export SERVER_ADDRESS=:8000

go run ./cmd/app
```

## Lessons

### Description

This project follows Go community standards for project organization. Learn more about recommended project layouts:

- [Go Project Layout](https://github.com/golang-standards/project-layout) — Standard Go project directory layout and
  conventions

### Alternative Database Libraries

Consider these Go database libraries for different use cases:

- [sqlx](https://github.com/jmoiron/sqlx) — Extensions to Go's standard `database/sql` library
- [Squirrel](https://github.com/Masterminds/squirrel) — SQL query builder with fluent API
- [Bun](https://bun.uptrace.dev/guide/query-update.html#api) — Fast and simple ORM for PostgreSQL, MySQL, and SQLite
- [GORM](https://gorm.io/docs/create.html) — Full-featured ORM with associations, hooks, and migrations

### Lesson 1

**Objective**: Extend the `UserStore` functionality to support additional user operations and address management.

#### Homework

**Task 1: Implement missing UserStore methods**

- GetUserByEmail
- ListUsersByEmail
- UpdateUser

**Task 2: Save address to the database**

Add functionality to store users' country and city in the database. **Optionally**, implement logic for
storing complete living addresses (country, city, street, zip code).

### Lesson 2

- [Database normalization](https://en.wikipedia.org/wiki/Database_normalization)
- [Pattern: Event sourcing](https://microservices.io/patterns/data/event-sourcing.html)
- [Implementing event sourcing using a relational database](https://softwaremill.com/implementing-event-sourcing-using-a-relational-database/)

#### Homework

**Task 1: Implement missing OperationStore methods**

- ListOperationsByUserID
- GetUserBalance
- CreateTransfer

### Lesson 3

[//]: # (Create job to calculate statistic. Show different ways to do that.)
