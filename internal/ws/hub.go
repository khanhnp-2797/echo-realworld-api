package ws

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const channelPrefix = "ws:comments:"

// Hub manages all active WebSocket rooms keyed by article slug.
// Broadcasts are published to Redis so all Hub instances (across multiple
// server replicas) deliver the same message to their local clients.
//
// Life-cycle:
//
//	hub := ws.NewHub(rdb)
//	go hub.Run()               // start the event loop (once, at startup)
//	hub.Broadcast(slug, msg)   // called by AddComment handler after DB write
type Hub struct {
	rdb *redis.Client

	mu      sync.RWMutex
	rooms   map[string]map[*Client]struct{} // slug → local clients
	pubsubs map[string]*redis.PubSub        // slug → active Redis subscription

	register   chan subscription
	unregister chan subscription
	deliver    chan roomMessage // inbound from Redis → local clients
}

type subscription struct {
	slug   string
	client *Client
}

type roomMessage struct {
	slug    string
	payload []byte
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		rdb:        rdb,
		rooms:      make(map[string]map[*Client]struct{}),
		pubsubs:    make(map[string]*redis.PubSub),
		register:   make(chan subscription, 64),
		unregister: make(chan subscription, 64),
		deliver:    make(chan roomMessage, 256),
	}
}

func (h *Hub) channel(slug string) string {
	return fmt.Sprintf("%s%s", channelPrefix, slug)
}

// Run processes register / unregister / deliver events sequentially.
// Must be started in a goroutine: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[s.slug]; !ok {
				h.rooms[s.slug] = make(map[*Client]struct{})
				// First local client for this slug — subscribe to Redis channel.
				ps := h.rdb.Subscribe(context.Background(), h.channel(s.slug))
				h.pubsubs[s.slug] = ps
				go h.listenRedis(s.slug, ps)
				log.Printf("[ws] redis subscribed channel=%s", h.channel(s.slug))
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
					// Last local client left — stop the Redis subscription.
					if ps, ok := h.pubsubs[s.slug]; ok {
						_ = ps.Close()
						delete(h.pubsubs, s.slug)
						log.Printf("[ws] redis unsubscribed channel=%s", h.channel(s.slug))
					}
				}
			}
			h.mu.Unlock()
			log.Printf("[ws] client unregistered slug=%s", s.slug)

		case rm := <-h.deliver:
			h.mu.RLock()
			room := h.rooms[rm.slug]
			h.mu.RUnlock()
			log.Printf("[ws] delivering slug=%s local_clients=%d", rm.slug, len(room))
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

// listenRedis reads messages from a Redis PubSub channel and forwards them
// to the deliver channel so Run() can distribute to local clients.
// Exits automatically when the PubSub is closed via ps.Close().
func (h *Hub) listenRedis(slug string, ps *redis.PubSub) {
	log.Printf("[ws] redis listener started slug=%s", slug)
	for msg := range ps.Channel() {
		log.Printf("[ws] redis message received channel=%s bytes=%d", msg.Channel, len(msg.Payload))
		h.deliver <- roomMessage{slug: slug, payload: []byte(msg.Payload)}
	}
	log.Printf("[ws] redis listener exiting slug=%s", slug)
}

// Broadcast publishes a JSON payload to the Redis Pub/Sub channel for the
// given article slug. Every Hub instance subscribed to that channel will
// deliver the message to its local clients.
func (h *Hub) Broadcast(slug string, payload []byte) {
	if err := h.rdb.Publish(context.Background(), h.channel(slug), payload).Err(); err != nil {
		log.Printf("[ws] redis publish error slug=%s: %v", slug, err)
		return
	}
	log.Printf("[ws] published to redis channel=%s bytes=%d", h.channel(slug), len(payload))
}

// Subscribe registers a client to a room. Safe to call from any goroutine.
func (h *Hub) Subscribe(slug string, c *Client) {
	h.register <- subscription{slug: slug, client: c}
}

// Unsubscribe removes a client from a room. Safe to call from any goroutine.
func (h *Hub) Unsubscribe(slug string, c *Client) {
	h.unregister <- subscription{slug: slug, client: c}
}
