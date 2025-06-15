package ws

import "github.com/gorilla/websocket"

// Client represents a websocket connection to a spreadsheet.
type Client struct {
	Conn    *websocket.Conn
	SheetID string
	Send    chan []byte
}