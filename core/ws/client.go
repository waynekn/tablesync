package ws

import (
	"log/slog"
	"strings"
	"sync"
	"time"

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

// NewClient instantiates and returns a new Client
func NewClient(sheetID string, colNum int, conn *websocket.Conn, collabStore *collab.Store, hub *Hub) *Client {
	client := &Client{
		Conn:        conn,
		SheetID:     sheetID,
		Send:        make(chan collab.EditMsg, 50),
		collabStore: collabStore,
		hub:         hub,
		done:        make(chan struct{}),
		closeOnce:   sync.Once{},
	}

	go client.readEdits()
	go client.writeEdits(colNum)
	return client
}

// readEdits listens for incoming edits from the client.
// It reads JSON messages from the websocket connection, applies them to Redis,
// and broadcasts them to other clients connected to the same sheet.
func (c *Client) readEdits() {
	defer func() {
		c.hub.Unregister <- c
		c.Close("")
	}()

	for {
		var edit collab.EditMsg
		err := c.Conn.ReadJSON(&edit)

		if websocket.IsCloseError(err) {
			break
		}

		if err != nil {
			slog.Error("failed to read JSON", "err", err)
			c.Close("The server was unable to read your edits")
			break
		}

		// Add 1 to the row index to account for the offset caused by how data is handled:
		// The server stores both the column headers and the sheet data in a single 2D array,
		// with headers at index 0. However, the client separates headers from data — it
		// shifts the first row to use as column headers. This causes a one-row difference
		// between client and server representations.
		redisEdit := collab.EditMsg{
			Row:  edit.Row + 1,
			Col:  edit.Col,
			Data: edit.Data,
		}

		err = c.collabStore.ApplyEdit(c.SheetID, redisEdit)
		if err != nil {
			slog.Error("error applying edit", "err", err)
			c.Close("Your changes couldn’t be saved due to a server error")
			break
		}
		c.hub.Broadcast <- collab.BroadCastMsg{SheetID: c.SheetID, Edit: edit}
	}
}

// writeEdits sends the initial sheet data to the client and listens for edits broadcasted
// to the clients `Send` channel and sends them to the client.
func (c *Client) writeEdits(colNum int) {
	defer func() {
		c.hub.Unregister <- c
		c.Close("")
	}()

	// collect previous sheet data in redis and send it to the client
	redisData, err := c.collabStore.GetRedisSheetData(c.SheetID)
	if err != nil {
		slog.Error("failed to retrieve sheet data", "sheetID", c.SheetID, "err", err)
		c.Close("Could not retrieve sheet data.")
		return
	}

	sheetData, err := mapToMatrix(redisData, colNum)
	if err != nil {
		c.Close("Unable to initialize sheet data")
		return
	}

	err = c.Conn.WriteJSON(sheetData)
	if err != nil {
		slog.Error("failed to send initial sheet data",
			"sheetID", c.SheetID,
			"err", err,
		)
		c.Close("The server was unable to send the initial sheet data")
		return
	}

	for {
		select {
		case edit := <-c.Send:
			err := c.Conn.WriteJSON(edit)
			if err != nil {
				c.Close("We couldn't send changes to your browser, please refresh to reconnect.")
				return
			}
		case <-c.done:
			return
		}
	}
}

// Close gracefully closes the websocket connection.
// It sends a close message with an optional reason and ensures that the connection is closed only once
func (c *Client) Close(reason string) {
	c.closeOnce.Do(func() {
		close(c.done)

		if strings.TrimSpace(reason) != "" {
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
			err := c.Conn.WriteControl(websocket.CloseMessage, cm, time.Now().Add(time.Second))
			if err != nil {
				slog.Error("failed to send close message", "err", err)
			}
		}

		c.Conn.Close()
	})
}
