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
GOBIN="$(pwd)/bin" go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5

PATH="$(pwd)/bin:${PATH}" go generate ./...
```
