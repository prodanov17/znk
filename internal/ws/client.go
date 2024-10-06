package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Convert the message to JSON and send it via WebSocket
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{}
		err = json.Unmarshal(m, msg)
		if err != nil {
			log.Printf("error unmarshaling: %v", err)
			continue
		}

		if err := handleMessage(hub, msg); err != nil {
			errorPayload, _ := json.Marshal(map[string]string{"error": err.Error()})
			hub.SendMessageToClient(&Message{Action: "error", Payload: errorPayload, UserID: c.ID})
			log.Printf("error handling message: %v", err)
			continue
		}

		fmt.Println("message", msg)
	}
}

func handleMessage(hub *Hub, msg *Message) error {
	action, err := CreateAction(msg.Action, msg.RoomID, msg.UserID, msg.Payload)
	if err != nil {
		return fmt.Errorf("error creating action: %v", err)
	}

	err = json.Unmarshal(msg.Payload, action)
	if err != nil {
		return fmt.Errorf("error unmarshalling payload: %v", err)
	}

	return action.Execute(hub)
}
