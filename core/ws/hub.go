package ws

// NewHub creates and returns a new Hub instance.
// The Hub manages websocket clients, allowing them to register, unregister,
// and broadcast messages to all clients connected to the same sheetID.
func NewHub() Hub {
	return Hub{
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client, 100),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan []byte, 100),
	}
}
