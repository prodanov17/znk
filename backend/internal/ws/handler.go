package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prodanov17/znk/internal/utils"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/ws/join/{lobbyID}", h.handleJoinLobby)
	router.HandleFunc("/ws/rooms", h.handleGetRooms)
	router.HandleFunc("/ws/clients", h.handleGetClients)
	router.HandleFunc("POST /ws/rooms/{room_id}/clear", h.handleClearRoom)
}

func (h *Handler) handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	//authenticate user

	lobbyID := r.PathValue("lobbyID")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	if clientID == "" || username == "" {
		utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("missing userId or username"))
		return
	}

	if lobbyID == "" {
		utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("missing or invalid lobbyID"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		Username: username,
		RoomID:   lobbyID,
	}
	fmt.Println("Registering client", cl.ID, cl.Username, cl.RoomID)

	h.hub.RegisterClient(cl)

	go cl.WriteMessage(h.hub)
	cl.ReadMessage(h.hub)
}

func (h *Handler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	//authenticate user

	rooms := make([]*Room, 0)

	for _, room := range h.hub.Room {
		rooms = append(rooms, room)
	}

	utils.WriteJSON(w, http.StatusOK, rooms)
}

func (h *Handler) handleGetClients(w http.ResponseWriter, _ *http.Request) {
	//authenticate user

	clients := make([]*Client, 0)

	for _, client := range h.hub.Clients {
		clients = append(clients, client)
	}

	utils.WriteJSON(w, http.StatusOK, clients)
}

func (h *Handler) handleClearRoom(w http.ResponseWriter, r *http.Request) {
	//authenticate user
	roomID := r.PathValue("room_id")

	err := h.hub.ClearRoom(roomID)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Rooms cleared")
}
