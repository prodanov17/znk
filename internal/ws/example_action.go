package ws

import (
	"encoding/json"
	"fmt"
)

type ExampleAction struct {
	BaseAction
	Payload string
}

func (a *ExampleAction) Execute(hub *Hub) error {
	payload := map[string]interface{}{"message": "example action"}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	hub.BroadcastMessage(&Message{Action: "example_result", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID})
	hub.SendMessageToClient(&Message{Action: "example_result_client", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID})
	return nil
}
