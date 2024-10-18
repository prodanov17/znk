package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/internal/utils"
	"github.com/prodanov17/znk/pkg/logger"
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
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/ws/join/{roomID}", h.handleJoinRoom)
	router.HandleFunc("/ws/rooms", h.handleGetRooms)
	router.HandleFunc("/ws/rooms/{room_id}", h.handleGetRoom)
	router.HandleFunc("/ws/clients", h.handleGetClients)
	router.HandleFunc("POST /ws/rooms/{room_id}/clear", h.handleClearRoom)
	router.HandleFunc("POST /ws/rooms", h.handleCreateRoom)
}

func (h *Handler) handleGetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("room_id")
	room, err := h.hub.RoomService().GetRoomByID(roomID)
	if err != nil {
		utils.WriteError(w, r, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, room)
}

func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomPayload types.CreateRoomPayload

	err := utils.ParseJSON(r, &createRoomPayload)
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	room, err := h.hub.RoomService().CreateRoom(&createRoomPayload)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, room)
}

func (h *Handler) handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	//authenticate user

	roomID := r.PathValue("roomID")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	if clientID == "" || username == "" {
		utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("missing userId or username"))
		return
	}

	if roomID == "" {
		utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("missing or invalid lobbyID"))
		return
	}

	_, err := h.hub.RoomService().GetRoomByID(roomID)
	if err != nil {
		utils.WriteError(w, r, http.StatusNotFound, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *types.Message, 1000),
		ID:       clientID,
		Username: username,
		RoomID:   roomID,
	}

	logger.Log.Infof("Client %s joined room %s | IP Addr: %s", cl.ID, cl.RoomID, r.RemoteAddr)

	h.hub.RegisterClient(cl)

	go cl.WriteMessage(h.hub)
	cl.ReadMessage(h.hub)
}

func (h *Handler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	//authenticate user

	rooms := make([]*types.Room, 0)

	roomSlice, _ := h.hub.RoomService().GetRooms()

	for _, room := range roomSlice {
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

	err := h.hub.RoomService().ClearRoom(roomID)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	logger.Log.Infof("Room %s cleared by %s", roomID, r.RemoteAddr)

	utils.WriteJSON(w, http.StatusOK, "Rooms cleared")
}
