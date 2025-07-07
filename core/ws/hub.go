package ws

import (
	"log/slog"
	"runtime/debug"
	"slices"

	"github.com/waynekn/tablesync/core/collab"
)

// Hub is a long-running in-memory struct that keeps track of all
// websocket clients currently connected. It handles broadcasting a message
// to all clients on the same sheetID
type Hub struct {
	Clients    map[string][]*Client // map[sheetID]clients
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan collab.BroadCastMsg
}

// NewHub creates and returns a new Hub instance.
// The Hub manages websocket clients, allowing them to register, unregister,
// and broadcast messages to all clients connected to the same sheetID.
func NewHub() *Hub {
	hub := &Hub{
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client, 100),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan collab.BroadCastMsg, 100),
	}
	go hub.run()
	return hub
}

// run starts the Hub's event loop, which listens for client registration,
// unregistration, and broadcast messages.
func (h *Hub) run() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Hub run recovered from panic", "err", r, "stack", string(debug.Stack()))
				}
			}()

			select {
			case client := <-h.Register:
				h.Clients[client.SheetID] = append(h.Clients[client.SheetID], client)
			case client := <-h.Unregister:
				i := 0
				h.Clients[client.SheetID] = slices.DeleteFunc(h.Clients[client.SheetID], func(c *Client) bool {
					i++
					return c == client
				})

				// if the sheet only had one client and they have been unregistred, delete the empty key
				if i == 1 {
					delete(h.Clients, client.SheetID)
				}
			case broadcast := <-h.Broadcast:
				if clients, ok := h.Clients[broadcast.SheetID]; ok {
					for _, client := range clients {
						client.Send <- broadcast.Msg
					}
				}
			}
		}()
	}
}
