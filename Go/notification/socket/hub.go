package socket

import (
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
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
			clientId := client.ID
			for client := range h.clients {
				msg := []byte("some one join room (ID: " + clientId + ")")
				client.Send <- msg
			}

			h.clients[client] = true

		case client := <-h.unregister:
			clientId := client.ID
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			for client := range h.clients {
				msg := []byte("some one leave room (ID:" + clientId + ")")
				client.Send <- msg
			}
		case userMessage := <-h.broadcast:
			var data map[string][]byte
			json.Unmarshal(userMessage, &data)

			for client := range h.clients {
				//prevent self receive the message
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
