package router

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/handlers"
	"github.com/waynekn/tablesync/api/middleware"
	"github.com/waynekn/tablesync/core/collab"
	"github.com/waynekn/tablesync/core/ws"
)

// Router holds dependencies and the Gin engine
type Router struct {
	engine *gin.Engine
	db     *sql.DB
	redis  *redis.Client
}

// New creates a new router with dependencies
func New(db *sql.DB, redis *redis.Client) *Router {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	router := Router{
		engine: r,
		db:     db,
		redis:  redis,
	}
	router.setupMiddleware()
	router.registerRoutes()
	return &router
}

// SetupMiddleware configures CORS and other middleware
func (r *Router) setupMiddleware() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.engine.Use(cors.New(config))
}

// RegisterRoutes sets up all application routes
func (r *Router) registerRoutes() {
	collabStore := collab.NewStore(r.redis)
	hub := ws.NewHub()

	// Initialize repositories
	spreadsheetRepo := repo.NewSpreadsheetRepo(r.db)
	wsRepo := repo.NewWsRepo(r.db)

	// Initialize handlers
	spreadsheetHandler := handlers.NewSpreadsheetHandler(spreadsheetRepo)
	wsHandler := handlers.NewWsHandler(wsRepo, collabStore, hub)

	// Register routes
	r.registerSpreadsheetRoutes(spreadsheetHandler)
	r.registerWebSocketRoutes(wsHandler)
}

func (r *Router) registerSpreadsheetRoutes(h *handlers.SpreadsheetHandler) {
	r.engine.POST("spreadsheet/create/", middleware.RequireAuth(r.redis), h.CreateSpreadsheetHandler)
	r.engine.GET("spreadsheets/", middleware.RequireAuth(r.redis), h.GetOwnSpreadsheetsHandler)
}

func (r *Router) registerWebSocketRoutes(h *handlers.WsHandler) {
	r.engine.GET("ws/sheet/:sheetID/edit/", h.EditSessionHandler)
}

// Run starts the server
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
