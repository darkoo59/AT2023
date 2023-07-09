package socket

import (
	"encoding/json"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) GetClient(id string) *Client {
	for client := range h.clients {
		if client.ID == id {
			return client
		}
	}
	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case userMessage := <-h.broadcast:
			var data map[string][]byte
			json.Unmarshal(userMessage, &data)

			for client := range h.clients {
				if client.ID == string(data["id"]) {
					select {
					case client.Send <- data["message"]:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
