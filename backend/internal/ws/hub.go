package ws

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/pkg/logger"
)

type Hub struct {
	sync.Mutex
	roomService types.RoomService
	Clients     map[string]*Client // client_id: client
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan *types.Message
}

func NewHub(roomService types.RoomService) *Hub {
	return &Hub{
		roomService: roomService,
		Clients:     make(map[string]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *types.Message, 10),
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

func (h *Hub) RoomService() types.RoomService {
	return h.roomService
}

func (h *Hub) RegisterClient(client *Client) {
	h.Register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Unregister <- client
}

func (h *Hub) BroadcastMessage(message *types.Message) {
	h.Broadcast <- message
}

func (h *Hub) registerClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	_, ok := h.Clients[client.ID]
	if ok {
		return
	}
	if client.RoomID != "" {
		room, err := h.roomService.GetRoomByID(client.RoomID)
		if err != nil { //rework game creation into its seperate method
			// room, err = h.roomService.CreateRoom(&types.CreateRoomPayload{UserID: client.ID, Username: client.Username})
			return
		}

		if len(room.Players) == 4 {
			message := &types.Message{Action: "room_full", Payload: nil, RoomID: client.RoomID, UserID: client.ID}
			h.SendMessageToClient(message)
			client.Disconnect()
			return
		}

		err = h.roomService.AddUserToRoom(client.RoomID, client.ID, client.Username)
		if err != nil {
			fmt.Println("Error adding user to room")
			return
		}

		game, err := h.roomService.GameService().GetGameByID(client.RoomID)
		player, err := game.Player(client.ID)
		if err != nil { //not found
			h.roomService.GameService().AddUserToGame(client.RoomID, player.UserID, player.Username)
		}

		logger.Log.Infof("Player %s joined room %s", client.Username, client.RoomID)

		h.Clients[client.ID] = client
		fmt.Println("got this far")

		playerTeam := game.PlayerTeam(client.ID)
		if playerTeam == nil {
			fmt.Println("Player team not found")
			NotFound(h)
			client.Disconnect()
			return
		}

		var playerPayload = map[string]interface{}{"id": client.ID, "username": client.Username, "team_id": playerTeam.ID, "teams": game.GameTeam, "playerCount": len(game.Players())}
		rawPayload, _ := json.Marshal(playerPayload)

		message := &types.Message{Action: "player_joined", Payload: rawPayload, RoomID: client.RoomID, UserID: client.ID}
		h.BroadcastMessage(message)
		h.SendMessageToClient(message)
	} else {
		fmt.Println("Room not found")
		NotFound(h)
		client.Disconnect()
		return
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	if client.RoomID == "" {
		return
	}
	logger.Log.Infof("Player %s left room %s", client.Username, client.RoomID)

	room, err := h.roomService.GetRoomByID(client.RoomID)
	if err != nil {
		logger.Log.Errorf("Error getting room: %v", err)
		return
	}

	game, err := h.roomService.GameService().GetGameByID(client.RoomID)
	playerTeam := game.PlayerTeam(client.ID)
	if playerTeam == nil {
		NotFound(h)
		return
	}

	h.roomService.RemoveUserFromRoom(client.RoomID, client.ID)
	if game.State != "started" {
		game.RemovePlayer(client.ID)
	}
	var playerPayload = map[string]interface{}{"id": client.ID, "username": client.Username, "team_id": playerTeam.ID, "teams": game.GameTeam, "playerCount": len(game.Players())}
	rawPayload, _ := json.Marshal(playerPayload)

	h.Broadcast <- &types.Message{Action: "player_left", Payload: rawPayload, RoomID: client.RoomID, UserID: client.ID}
	delete(h.Clients, client.ID)
	// Check if room is empty and set a 5-second timeout to delete the room if empty
	if len(room.Players) == 0 {
		go h.scheduleRoomDeletion(client.RoomID, 5*time.Second)
	}
}

// scheduleRoomDeletion waits for a specified duration and then deletes the room if it's still empty
func (h *Hub) scheduleRoomDeletion(roomID string, delay time.Duration) {
	// Wait for the specified duration
	time.Sleep(delay)

	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// Re-check if the room is still empty
	room, err := h.roomService.GetRoomByID(roomID)
	if err != nil {
		logger.Log.Errorf("Error getting room: %v", err)
		return
	}

	if len(room.Players) == 0 {
		// Room is empty, proceed to delete
		logger.Log.Infof("Room %s is empty. Deleting room...", roomID)
		h.roomService.ClearRoom(roomID)
	} else {
		logger.Log.Infof("Room %s is not empty. Skipping room deletion.", roomID)
	}
}

func (h *Hub) broadcastMessage(message *types.Message) {
	if message.RoomID != "" {
		room, err := h.roomService.GetRoomByID(message.RoomID)
		if err != nil {
			logger.Log.Errorf("Error getting room: %v", err)
			return
		}

		for _, player := range room.Players {
			if player.UserID == message.UserID {
				continue
			}
			client, ok := h.Clients[player.UserID]
			if !ok {
				continue
			}
			client.Message <- message
		}
	} else {
		logger.Log.Info("Broadcasting message to all clients")
		for _, client := range h.Clients {
			if client.ID == message.UserID {
				continue
			}

			client.Message <- message
		}
	}
}

func (h *Hub) SendMessageToClient(message *types.Message) {
	client, ok := h.Clients[message.UserID]
	if !ok {
		return
	}

	client.Message <- message
}

func NotFound(h *Hub) {
	fmt.Println("Not found")
	errorPayload := map[string]interface{}{"error": " not found"}
	rawPayload, _ := json.Marshal(errorPayload)

	h.Broadcast <- &types.Message{Action: "error", Payload: rawPayload}
}
