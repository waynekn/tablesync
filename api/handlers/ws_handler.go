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
			closeWsConn("The sheet you're trying to edit does not exist.", conn)
		} else {
			closeWsConn("An unexpected error occurred while connecting. Please try again later.", conn)
		}
		return
	}

	now := time.Now().UTC()

	if now.After(sheet.Deadline) {
		closeWsConn("The deadline to edit this sheet has passed.", conn)
		return
	}

	exists, err := h.collab.SheetExists(sheetID)

	if err != nil {
		closeWsConn("An error occurred during initialization. Please try again later.", conn)
		return
	}

	var sheetData [][]string
	err = json.Unmarshal(sheet.Data, &sheetData)
	if err != nil {
		slog.Error("error unmarshalling sheet data", "err", err)
		closeWsConn("Could not process sheet data. Please try again in a while", conn)
		return
	}

	if !exists {
		err = h.collab.InitRedisSheet(sheetID, sheet.Deadline, &sheetData)
		if err != nil {
			slog.Error("error initializing redis sheet", "err", err)
			closeWsConn("Could not initialize collaborative session.", conn)
			return
		}
	}

	cols := sheetData[0]
	client := ws.NewClient(sheetID, len(cols), conn, h.collab, h.hub)
	h.hub.Register <- client
}

// closeWsConn closes the WebSocket connection with the provided reason.
// This function is used to properly close an existing WebSocket connection
// before a new ws.Client instance is created, which has its own Close method.
// The connection is closed by sending a close message with the provided reason
// and a controlled close timeout.
func closeWsConn(reason string, conn *websocket.Conn) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	err := conn.WriteControl(websocket.CloseMessage, cm, time.Now().Add(time.Second))
	if err != nil {
		slog.Error("failed to send close message", "err", err)
	}
	conn.Close()
}
