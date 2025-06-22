package ws

import (
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	hub := NewHub()

	client := &Client{
		Conn:    &websocket.Conn{},
		SheetID: "test-sheet",
		Send:    make(chan []byte, 1),
	}

	var wg sync.WaitGroup

	// Start a goroutine to wait for the client to be registered
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Register the client
		hub.Register <- client
	}()

	// Wait for the register operation to be processed
	time.Sleep(50 * time.Millisecond)

	// Check if the client is registered
	assert.Equal(t, 1, len(hub.Clients[client.SheetID]),
		"After registering the client, the client count for SheetID '%s' should be 1", client.SheetID)

	// unregister the client
	wg.Add(1)
	go func() {
		defer wg.Done()
		hub.Unregister <- client
	}()

	// Wait for the unregister operation
	time.Sleep(50 * time.Millisecond)

	// Check if the client is unregistered
	assert.Equal(t, 0, len(hub.Clients[client.SheetID]),
		"After unregistering the client, the client count for SheetID '%s' should be 0", client.SheetID)

	// check that the now sheet has been removed
	_, exists := hub.Clients[client.SheetID]
	assert.False(t, exists, "Expected sheetID %q to be removed after last client unregistered", client.SheetID)

	wg.Wait()
}
