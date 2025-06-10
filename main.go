package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/waynekn/tablesync/api"
	"github.com/waynekn/tablesync/api/db"
	"github.com/waynekn/tablesync/api/logging"
	"github.com/waynekn/tablesync/api/router"
	"github.com/waynekn/tablesync/core/rdb"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		slog.Error("Failed to convert REDIS_DB to int", "err", err)
		os.Exit(1)
	}

	redisClient, err := rdb.Connect(redisAddr, redisPassword, redisDB)
	if err != nil {
		slog.Error("Redis connection failed, shutting down", "addr", redisAddr, "err", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	api.RegisterJSONTagNameFormatter()

	router := router.New(conn, redisClient)

	router.Run("localhost:8000")
}
