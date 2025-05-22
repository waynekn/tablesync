package router

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/handlers"
	"github.com/waynekn/tablesync/api/middleware"
)

// New registers the routes and returns the router.
func New(db *sql.DB) *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// Set up the routes
	// Initialize repositories with database connection
	spreadsheetRepo := repo.NewSpreadsheetRepo(db)
	spreadsheetHandler := handlers.NewSpreadsheetHandler(spreadsheetRepo)

	r.POST("spreadsheet/create/", middleware.RequireAuth(), spreadsheetHandler.CreateSpreadsheetHandler)

	r.SetTrustedProxies([]string{"127.0.0.1"})

	return r
}
