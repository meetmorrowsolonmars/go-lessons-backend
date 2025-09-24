//go:build integration

package postgres

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbDatabase = "test"
	dbUsername = "user"
	dbPassword = "password"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func TestMain(m *testing.M) {
	// Setup containers.
	log.Println("Setup containers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine3.22",
		postgres.WithDatabase(dbDatabase),
		postgres.WithUsername(dbUsername),
		postgres.WithPassword(dbPassword),
		testcontainers.WithCmdArgs("-c", "log_statement=all"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalln("Run container:", err)
	}

	// Create a database connection.
	log.Println("Create a database connection...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalln("Get connection string:", err)
	}

	dbConnection, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalln("Create connection:", err)
	}

	once.Do(func() {
		db = dbConnection
	})

	// Migrate.
	{
		log.Println("Migrate...")

		db := stdlib.OpenDBFromPool(dbConnection)
		defer func() {
			_ = db.Close()
		}()

		dir, _ := os.Getwd()
		log.Println("Current directory:", dir)

		err = goose.UpContext(ctx, db, dir+"/../../../migrations")
		if err != nil {
			log.Fatalln("Migrate:", err)
		}
	}

	// Run tests.
	log.Println("Database connection string:", connString)
	log.Println("Run tests...")

	code := m.Run()

	// Terminate containers.
	log.Println("Terminate containers...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = postgresContainer.Terminate(ctx)
	if err != nil {
		log.Fatalln("Terminate container:", err)
	}

	// Exit.
	os.Exit(code)
}
