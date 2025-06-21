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

// EditMsg carries the details of a spreadsheet cell edit
// made by a client, for broadcast to other collaborators.
type EditMsg struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Data string `json:"data"`
}
