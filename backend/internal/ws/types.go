package ws

import (
	"encoding/json"
)

type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
	UserID  string          `json:"user_id"`
	RoomID  string          `json:"room_id"`
}

type Response struct {
	Action  string          `json:"action"`
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
}

type ActionPayload struct {
	Action  string `json:"action"`
	UserID  int    `json:"user_id"`
	RoomID  string `json:"room_id"`
	LobbyID string `json:"lobby_id"`
}
