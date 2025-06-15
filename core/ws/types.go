package ws

import "github.com/gorilla/websocket"

// Client represents a websocket connection to a spreadsheet.
type Client struct {
	Conn    *websocket.Conn
	SheetID string
	Send    chan []byte
}

// Hub is a long-running in-memory struct that keeps track of all
// websocket clients currently connected. It handles broadcasting a message
// to all clients on the same sheetID
type Hub struct {
	Clients    map[string][]*Client // map[sheetID]clients
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}
