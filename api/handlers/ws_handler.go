package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/models"
	"github.com/waynekn/tablesync/core/collab"
	"github.com/waynekn/tablesync/core/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:5173"
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WsHandler struct {
	repo   repo.WsRepo
	collab *collab.Store
	hub    *ws.Hub
}

// NewWsHandler creates a new instance of WsHandler with the provided repository and collaboration store.
func NewWsHandler(repo repo.WsRepo, collabStore *collab.Store, hub *ws.Hub) *WsHandler {
	return &WsHandler{repo: repo, collab: collabStore, hub: hub}
}

// EditSessionHandler initializes WebSocket connections for editing a spreadsheet.
// It upgrades the HTTP connection to a WebSocket connection and checks if the
// specified spreadsheet exists and is still editable (i.e., the deadline has not passed).
func (h *WsHandler) EditSessionHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Error upgrading connection", "err", err)
		return
	}

	sheetID := c.Param("sheetID")
	sheet, err := h.repo.GetSheetByID(sheetID)

	if err != nil {
		if err == sql.ErrNoRows {
			conn.WriteJSON(models.NewWsErrMsg("The sheet you're trying to edit does not exist."))
		} else {
			conn.WriteJSON(models.NewWsErrMsg("An unexpected error occurred while connecting. Please try again later."))
		}
		conn.Close()
		return
	}

	now := time.Now().UTC()

	if now.After(sheet.Deadline) {
		conn.WriteJSON(models.NewWsErrMsg("The deadline to edit this sheet has passed."))
		conn.Close()
		return
	}

	exists, err := h.collab.SheetExists(sheetID)

	if err != nil {
		conn.WriteJSON(models.NewWsErrMsg("An error occurred during initialization. Please try again later."))
		conn.Close()
		return
	}

	if !exists {
		var sheetData [][]string
		err = json.Unmarshal(sheet.Data, &sheetData)
		if err != nil {
			slog.Error("error unmarshalling sheet data", "err", err)
			conn.WriteJSON(models.NewWsErrMsg("Could not process sheet data. Please try again in a while"))
			conn.Close()
			return
		}

		err = h.collab.InitRedisSheet(sheetID, sheet.Deadline, &sheetData)
		if err != nil {
			slog.Error("error initializing redis sheet", "err", err)
			conn.WriteJSON(models.NewWsErrMsg("Could not initialize collaborative session."))
			conn.Close()
			return
		}
	}

	client := ws.NewClient(sheetID, conn, h.collab, h.hub)
	h.hub.Register <- client
}
