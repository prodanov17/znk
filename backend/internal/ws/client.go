package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn `json:"-"`
	Message  chan *Message   `json:"-"`
	ID       string          `json:"id"`
	Username string          `json:"username"`
	RoomID   string          `json:"room_id"`
}

func (c *Client) WriteMessage(hub *Hub) {
	defer func() {
		c.Disconnect()
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
				hub.UnregisterClient(c)

				return
			}
		}
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.UnregisterClient(c)
		c.Disconnect()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println("Received message: ", string(m[:]))

		msg := &Message{}
		err = json.Unmarshal(m, msg)
		if err != nil {
			log.Printf("error unmarshaling: %v", err)
			continue
		}

		log.Printf("received message: %v", msg.Action)

		if err := handleMessage(hub, msg); err != nil {
			errorPayload, _ := json.Marshal(map[string]string{"error": err.Error()})
			hub.SendMessageToClient(&Message{Action: "error", Payload: errorPayload, UserID: c.ID})
			log.Printf("error handling message: %v", err)
			continue
		}

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

func (c *Client) Disconnect() {
	c.Conn.Close()
}
