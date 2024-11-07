package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type HandlerGroup struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

func NewHandlerGroup() *HandlerGroup {
	return &HandlerGroup{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

func (g *HandlerGroup) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", g.handleConnections)

	return mux
}
