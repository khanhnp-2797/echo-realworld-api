package ws

import (
	"log"
	"sync"
)

// Hub manages all active WebSocket rooms keyed by article slug.
// Each room holds the set of clients currently watching that article's comments.
//
// Life-cycle:
//
//	hub := ws.NewHub()
//	go hub.Run()               // start the event loop (once, at startup)
//	hub.Broadcast(slug, msg)   // called by AddComment handler after DB write
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{} // slug → set of clients

	register   chan subscription
	unregister chan subscription
	broadcast  chan roomMessage
}

type subscription struct {
	slug   string
	client *Client
}

type roomMessage struct {
	slug    string
	payload []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]struct{}),
		register:   make(chan subscription, 64),
		unregister: make(chan subscription, 64),
		broadcast:  make(chan roomMessage, 256),
	}
}

// Run processes register / unregister / broadcast events sequentially.
// Must be started in a goroutine: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[s.slug]; !ok {
				h.rooms[s.slug] = make(map[*Client]struct{})
			}
			h.rooms[s.slug][s.client] = struct{}{}
			h.mu.Unlock()
			log.Printf("[ws] client registered slug=%s total_in_room=%d", s.slug, len(h.rooms[s.slug]))

		case s := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[s.slug]; ok {
				delete(room, s.client)
				if len(room) == 0 {
					delete(h.rooms, s.slug)
				}
			}
			h.mu.Unlock()
			log.Printf("[ws] client unregistered slug=%s", s.slug)

		case rm := <-h.broadcast:
			h.mu.RLock()
			room := h.rooms[rm.slug]
			h.mu.RUnlock()
			log.Printf("[ws] broadcasting slug=%s clients=%d", rm.slug, len(room))
			for c := range room {
				select {
				case c.send <- rm.payload:
					log.Printf("[ws] queued message for client slug=%s", rm.slug)
				default:
					// Client channel full; close it so the write pump cleans up.
					log.Printf("[ws] client send buffer full, dropping slug=%s", rm.slug)
					close(c.send)
				}
			}
		}
	}
}

// Broadcast pushes a JSON payload to every client watching the given article slug.
func (h *Hub) Broadcast(slug string, payload []byte) {
	h.broadcast <- roomMessage{slug: slug, payload: payload}
}

// Subscribe registers a client to a room. Safe to call from any goroutine.
func (h *Hub) Subscribe(slug string, c *Client) {
	h.register <- subscription{slug: slug, client: c}
}

// Unsubscribe removes a client from a room. Safe to call from any goroutine.
func (h *Hub) Unsubscribe(slug string, c *Client) {
	h.unregister <- subscription{slug: slug, client: c}
}
