package ws

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/waynekn/tablesync/core/collab"
)

// Client represents a websocket connection to a spreadsheet.
type Client struct {
	Conn        *websocket.Conn
	SheetID     string
	Send        chan collab.EditMsg
	collabStore *collab.Store
	hub         *Hub
	done        chan struct{}
	closeOnce   sync.Once
}

// NewClient instantiates and returns a new Client
func NewClient(sheetID string, conn *websocket.Conn, collabStore *collab.Store, hub *Hub) *Client {
	return &Client{
		Conn:        conn,
		SheetID:     sheetID,
		Send:        make(chan collab.EditMsg, 50),
		collabStore: collabStore,
		hub:         hub,
		done:        make(chan struct{}),
		closeOnce:   sync.Once{},
	}
}

// Close gracefully closes the websocket connection.
// It sends a close message with an optional reason and ensures that the connection is closed only once
func (c *Client) Close(reason string) {
	c.closeOnce.Do(func() {
		close(c.done)

		if strings.TrimSpace(reason) != "" {
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
			err := c.Conn.WriteControl(websocket.CloseMessage, cm, time.Now().Add(time.Second))
			if err != nil {
				slog.Error("failed to send close message", "err", err)
			}
		}

		c.Conn.Close()
	})
}
