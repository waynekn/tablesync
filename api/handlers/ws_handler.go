package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/waynekn/tablesync/api/db/repo"
	"github.com/waynekn/tablesync/api/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:5173"
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WsHandler struct {
	repo repo.WsRepo
}

// NewWsHandler creates a new instance of WsHandler with the provided repository.
func NewWsHandler(repo repo.WsRepo) *WsHandler {
	return &WsHandler{repo: repo}
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
}
