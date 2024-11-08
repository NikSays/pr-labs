package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func (g *HandlerGroup) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	g.clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if msg.Message == "/leave" || err != nil {
			delete(g.clients, conn)
			return
		}

		g.broadcast <- msg
	}
}

func (g *HandlerGroup) HandleMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-g.broadcast:
			for client := range g.clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Print(err)
					client.Close()
					delete(g.clients, client)
				}
			}
		}
	}
}
