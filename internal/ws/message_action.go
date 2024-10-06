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
	Message string `json:"message"`
}

func (a *MessageAction) Execute(hub *Hub) error {
	fmt.Println("Received payload:", string(a.Payload))
	var messagePayload MessagePayload // Create a variable to hold the unmarshalled payload

	err := json.Unmarshal(a.Payload, &messagePayload) // Unmarshal into the variable
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	fmt.Println("Unmarshalled payload:", messagePayload)
	rawPayload, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	hub.BroadcastMessage(&Message{Action: "message", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID})
	return nil
}
