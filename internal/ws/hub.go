package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/prodanov17/znk/internal/game"
)

type Room struct {
	*game.Game
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
	sync.Mutex
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
	fmt.Println("Broadcasting message", message.Action)
	h.Broadcast <- message
}

func (h *Hub) registerClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	log.Println("ID", client.ID)
	_, ok := h.Clients[client.ID]
	if ok {
		log.Println("Client already registered")
		return
	}
	if client.RoomID != "" {
		room, ok := h.Room[client.RoomID]
		if !ok {
			h.Room[client.RoomID] = &Room{RoomID: client.RoomID, Game: game.NewGame(client.RoomID, client.ID), Clients: []*Client{}}
			room = h.Room[client.RoomID]
		}

		if len(room.Clients) == 4 {
			log.Println("Room is full")
			message := &Message{Action: "room_full", Payload: nil, RoomID: client.RoomID, UserID: client.ID}
			h.SendMessageToClient(message)
			client.Disconnect()
			return
		}

		room.AddClient(client)

		player, err := room.Game.Player(client.ID)
		if err != nil { //not found
			player = &game.Player{UserID: client.ID, Username: client.Username}
			room.Game.AddPlayer(player)
		}
		fmt.Println("Player added to game", room.RoomID)

		h.Clients[client.ID] = client

		playerTeam := room.Game.PlayerTeam(client.ID)
		if playerTeam == nil {
			NotFound(h)
			client.Disconnect()
			return
		}

		var playerPayload = map[string]interface{}{"id": client.ID, "username": client.Username, "team_id": playerTeam.ID, "teams": room.Game.GameTeam, "playerCount": len(room.Game.Players())}
		rawPayload, _ := json.Marshal(playerPayload)

		message := &Message{Action: "player_joined", Payload: rawPayload, RoomID: client.RoomID, UserID: client.ID}
		fmt.Println("Registering client", h.Clients)
		h.BroadcastMessage(message)
		h.SendMessageToClient(message)
	} else {
		NotFound(h)
		log.Println("Room ID is required")
		client.Disconnect()
		return
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	if client.RoomID == "" {
		log.Println("Room ID is required")
		return
	}
	fmt.Println("Unregistering client", client.ID)

	room, ok := h.Room[client.RoomID]
	if !ok {
		NotFound(h)
		return
	}

	playerTeam := room.Game.PlayerTeam(client.ID)
	if playerTeam == nil {
		NotFound(h)
		return
	}

	room.RemoveClient(client)
	if room.Game.State != "started" {
		room.Game.RemovePlayer(client.ID)
	}
	var playerPayload = map[string]interface{}{"id": client.ID, "username": client.Username, "team_id": playerTeam.ID, "teams": room.Game.GameTeam, "playerCount": len(room.Game.Players())}
	rawPayload, _ := json.Marshal(playerPayload)

	h.Broadcast <- &Message{Action: "player_left", Payload: rawPayload, RoomID: client.RoomID, UserID: client.ID}
	delete(h.Clients, client.ID)
	if len(room.Clients) == 0 {
		delete(h.Room, client.RoomID)
	}

}

func (h *Hub) broadcastMessage(message *Message) {
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
