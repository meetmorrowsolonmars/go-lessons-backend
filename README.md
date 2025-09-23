# Go Lessons Backend

This repository contains a series of workshops designed to teach Go backend development through a hands-on pet project:
a simple expense tracking application. Each workshop builds upon the previous one, gradually introducing more advanced
topics and technologies.

## Getting Started

Clone the repository:

```shell
git clone https://github.com/meetmorrowsolonmars/go-lessons-backend.git
cd go-lessons-backend
```

Navigate to the desired workshop branch (e.g., `workshop-1`):

```shell
git checkout workshop-1
```

Build and run the server:

```shell
go build -o ./bin/server ./cmd/server/main.go
APP_PUBLIC_ADDRESS=:6000 ./bin/server
```

```shell
openssl rand -base64 256 > ./.jwt_secret_key
```

Install tools

```shell
GOBIN="$(pwd)/bin" go install github.com/pressly/goose/v3/cmd/goose@v3.24.3
GOBIN="$(pwd)/bin" go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5

PATH="$(pwd)/bin:${PATH}" goose -dir=migrations create create_user_table sql
PATH="$(pwd)/bin:${PATH}" go generate ./...
```

```shell
export GOOSE_DBSTRING="dbname=postgres user=postgres password=postgres host=127.0.0.1 port=25432 sslmode=disable"

PATH="$(pwd)/bin:${PATH}" goose -dir=migrations postgres status
PATH="$(pwd)/bin:${PATH}" goose -dir=migrations postgres up
```
