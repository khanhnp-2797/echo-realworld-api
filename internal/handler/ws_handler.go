package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	appws "github.com/khanhnp-2797/echo-realworld-api/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins in development.
	// In production, restrict to your front-end origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler handles WebSocket upgrade requests.
type WSHandler struct {
	hub *appws.Hub
}

func NewWSHandler(hub *appws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// ServeComments upgrades the connection and subscribes the client to the
// comment room for the given article slug.
//
// GET /ws/articles/:slug/comments
//
// Once connected the client receives real-time JSON events whenever a new
// comment is posted to that article:
//
//	{
//	  "type": "new_comment",
//	  "comment": { ...CommentBody... }
//	}
func (h *WSHandler) ServeComments(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug is required")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		// upgrader already wrote the HTTP error response
		return nil
	}

	client := appws.NewClient(h.hub, slug, conn)
	h.hub.Subscribe(slug, client)

	// writePump owns the connection from here; readPump handles close/ping.
	go client.WritePump()
	go client.ReadPump()

	return nil
}
