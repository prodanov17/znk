package ws

import (
	"encoding/json"
	"fmt"
)

type MessageAction struct {
	BaseAction
	Payload json.RawMessage
}

type MessagePayload struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

func (a *MessageAction) Execute(hub *Hub) error {
	var messagePayload MessagePayload // Create a variable to hold the unmarshalled payload

	err := json.Unmarshal(a.Payload, &messagePayload) // Unmarshal into the variable
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	username := hub.Clients[a.UserID].Username
	messagePayload.Username = username

	if messagePayload.Message == "" {
		return fmt.Errorf("message is empty")
	}
	rawPayload, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	message := &Message{Action: "message", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	hub.BroadcastMessage(message)
	hub.SendMessageToClient(message)
	return nil
}
