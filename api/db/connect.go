package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)


var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database using the
// pgx driver. It retrieves the connection string from the DATABASE_URL
// environment variable, verifies the connection with a ping, and returns
// a *sql.DB instance.
//
// It returns an error if the environment variable is missing, the driver
// fails to open, or the database ping fails.
func Connect() (*sql.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		return nil, errors.New("missing DATABASE_URL")
	}

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		slog.Error("Failed to open pgx driver", "error", err)
		return nil, err
	}

	slog.Info("Pinging the database...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("Failed to ping the database", "error", err)
		return nil, err
	}

	slog.Info("Successfully connected to the database")

	DB = db
	return db, nil
}
