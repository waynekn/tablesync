package ws

import (
	"sync"

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

// Hub is a long-running in-memory struct that keeps track of all
// websocket clients currently connected. It handles broadcasting a message
// to all clients on the same sheetID
type Hub struct {
	Clients    map[string][]*Client // map[sheetID]clients
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan collab.BroadCastMsg
}
