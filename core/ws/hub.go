package ws

import (
	"slices"
)

// NewHub creates and returns a new Hub instance.
// The Hub manages websocket clients, allowing them to register, unregister,
// and broadcast messages to all clients connected to the same sheetID.
func NewHub() *Hub {
	hub := &Hub{
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client, 100),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan BroadCastMsg, 100),
	}
	go hub.run()
	return hub
}

// run starts the Hub's event loop, which listens for client registration,
// unregistration, and broadcast messages.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.SheetID] = append(h.Clients[client.SheetID], client)
		case client := <-h.Unregister:
			h.Clients[client.SheetID] = slices.DeleteFunc(h.Clients[client.SheetID], func(c *Client) bool {
				return c == client
			})
		}
	}
}
