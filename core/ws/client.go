package ws

import "github.com/gorilla/websocket"

// NewClient instantiates and returns a new Client
func NewClient(sheetID string, conn *websocket.Conn) *Client {
	return &Client{
		Conn:    conn,
		SheetID: sheetID,
		Send:    make(chan []byte, 50),
	}
}
