package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/prodanov17/znk/internal/game"
)

type Room struct {
	game.Game
	RoomID  string
	Clients []*Client
}

func (r *Room) AddClient(client *Client) {
	r.Clients = append(r.Clients, client)
}

func (r *Room) RemoveClient(client *Client) {
	for i, c := range r.Clients {
		if c == client {
			r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
			break
		}
	}
}

type Hub struct {
	Room       map[string]*Room   // _id:
	Clients    map[string]*Client // client_id: client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Room:       make(map[string]*Room),
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	// h.Room["1234"] = &Room{RoomID: "1234", Game:*game.NewGame(1, ), Clients: []*Client{}}
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.Register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Unregister <- client
}

func (h *Hub) BroadcastMessage(message *Message) {
	h.Broadcast <- message
}

func (h *Hub) registerClient(client *Client) {
	log.Println("ID", client.RoomID)
	if client.RoomID != "" {
		room, ok := h.Room[client.RoomID]
		if !ok {
			h.Room[client.RoomID] = &Room{RoomID: client.RoomID, Game: *game.NewGame(client.RoomID, client.ID), Clients: []*Client{}}
			room = h.Room[client.RoomID]
		}

		room.AddClient(client)

		player := game.Player{UserID: client.ID, Username: client.Username}
		room.Game.AddPlayer(&player)
		fmt.Println("Player added to game", room.RoomID)

	}

	h.Clients[client.ID] = client

	var playerPayload = map[string]interface{}{"id": client.ID, "username": client.Username}
	rawPayload, _ := json.Marshal(playerPayload)

	h.BroadcastMessage(&Message{Action: "player_joined", Payload: rawPayload, RoomID: client.RoomID, UserID: client.ID})
}

func (h *Hub) unregisterClient(client *Client) {
	if client.RoomID != "" {
		room, ok := h.Room[client.RoomID]
		if !ok {
			NotFound(h)
			return
		}
		room.RemoveClient(client)
		room.Game.RemovePlayer(client.ID)
	}

	delete(h.Clients, client.ID)
	h.Broadcast <- &Message{Action: "player_left", RoomID: client.RoomID, UserID: client.ID}
}

func (h *Hub) broadcastMessage(message *Message) {
	fmt.Println("Broadcasting message", message.Action)
	fmt.Println("Room", len(h.Room[message.RoomID].Clients))
	if message.RoomID != "" {
		room, ok := h.Room[message.RoomID]
		if !ok {
			NotFound(h)
			return
		}

		for _, client := range room.Clients {
			fmt.Println(client.ID, message.UserID)
			if client.ID == message.UserID {
				continue
			}
			client.Message <- message
		}
	} else {
		log.Println("Broadcasting message to all clients")
		for _, client := range h.Clients {
			if client.ID == message.UserID {
				continue
			}

			client.Message <- message
		}
	}
}

func (h *Hub) SendMessageToClient(message *Message) {
	client, ok := h.Clients[message.UserID]
	if !ok {
		return
	}

	client.Message <- message
}

func NotFound(h *Hub) {
	errorPayload := map[string]interface{}{"error": " not found"}
	rawPayload, _ := json.Marshal(errorPayload)

	h.Broadcast <- &Message{Action: "error", Payload: rawPayload}
}
