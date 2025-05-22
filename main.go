package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/waynekn/tablesync/api"
	"github.com/waynekn/tablesync/api/db"
	"github.com/waynekn/tablesync/api/logging"
	"github.com/waynekn/tablesync/api/router"
)

func main() {
	logging.InitLogger()

	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load the env vars", "error", err)
		os.Exit(1)
	}

	conn, err := db.Connect()
	if err != nil {
		slog.Error("Database connection failed, shutting down")
		os.Exit(1)
	}

	defer conn.Close()

	api.RegisterJSONTagNameFormatter()

	router := router.New(conn)

	router.Run("localhost:8000")
}
