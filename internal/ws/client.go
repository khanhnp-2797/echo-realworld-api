package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 // bytes — clients only send pong frames, not data
)

// Client represents a single WebSocket connection subscribed to one article room.
type Client struct {
	hub  *Hub
	slug string
	conn *websocket.Conn
	send chan []byte // buffered channel of outbound JSON messages
}

func NewClient(hub *Hub, slug string, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		slug: slug,
		conn: conn,
		send: make(chan []byte, 64),
	}
}

// ReadPump keeps the connection alive by consuming control frames (ping/pong/close).
// It unregisters the client when the connection closes.
func (c *Client) ReadPump() {
	defer func() {
		log.Printf("[ws] ReadPump exiting slug=%s — unsubscribing", c.slug)
		c.hub.Unsubscribe(c.slug, c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	// Drain incoming frames — clients are read-only in this design.
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("[ws] unexpected close slug=%s: %v", c.slug, err)
			} else {
				log.Printf("[ws] read error slug=%s: %v", c.slug, err)
			}
			break
		}
	}
}

// WritePump forwards buffered messages to the WebSocket and sends periodic pings.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("[ws] WritePump exiting slug=%s", c.slug)
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[ws] write error slug=%s: %v", c.slug, err)
				return
			}
			log.Printf("[ws] message sent slug=%s bytes=%d", c.slug, len(msg))

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
