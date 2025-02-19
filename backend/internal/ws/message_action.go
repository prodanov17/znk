package ws

import (
	"encoding/json"
	"fmt"

	"github.com/prodanov17/znk/internal/types"
)

type MessageAction struct {
	BaseAction
	Payload json.RawMessage
}

type MessagePayload struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

func (a *MessageAction) Execute() error {
	var messagePayload MessagePayload // Create a variable to hold the unmarshalled payload

	err := json.Unmarshal(a.Payload, &messagePayload) // Unmarshal into the variable
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	player, err := a.Hub.roomService.GetPlayerById(a.RoomID, a.UserID)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	username := player.Username
	messagePayload.Username = username

	if messagePayload.Message == "" {
		return fmt.Errorf("message is empty")
	}
	rawPayload, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	message := &types.Message{Action: "message", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	a.Hub.BroadcastMessage(message)
	a.Hub.SendMessageToClient(message)
	return nil
}
