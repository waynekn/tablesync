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
	Broadcast  chan BroadCastMsg
}

// EditMsg carries the details of a spreadsheet cell edit
// made by a client, for broadcast to other collaborators.
type EditMsg struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Data string `json:"data"`
}

// BroadCastMsg represents a message to be broadcasted to sheet collaborators.
// It has a SheetID field to identify the sheet whose clients should receive the
// message and a Msg field, which is an EditMsg, holding the data that will be
// broadcasted.
type BroadCastMsg struct {
	SheetID string
	Msg     EditMsg
}
