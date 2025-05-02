package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/waynekn/tablesync/api/logging"
	"github.com/waynekn/tablesync/api/router"
)

func main() {
	logging.InitLogger()

	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load the env vars", "error", err)
		os.Exit(1)
	}

	router := router.New()

	router.Run("localhost:8000")
}
