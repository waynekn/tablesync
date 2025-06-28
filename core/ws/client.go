package ws

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/waynekn/tablesync/core/collab"
)

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
