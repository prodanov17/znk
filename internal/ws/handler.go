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
	router.HandleFunc("/ws/join/{lobbyID}", h.JoinLobby)
}

func (h *Handler) JoinLobby(w http.ResponseWriter, r *http.Request) {
	//authenticate user

	fmt.Println("Joining lobby")

	lobbyID := r.PathValue("lobbyID")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	if clientID == "" || username == "" {
		utils.WriteError(w, r, http.StatusBadRequest, fmt.Errorf("missing userId or username"))
		return
	}

	if lobbyID == "" || lobbyID != "1234" {
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
	fmt.Println("Registering client")

	h.hub.RegisterClient(cl)

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)
}
