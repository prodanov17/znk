package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/pkg/logger"
)

type Client struct {
	Conn     *websocket.Conn     `json:"-"`
	Message  chan *types.Message `json:"-"`
	ID       string              `json:"id"`
	Username string              `json:"username"`
	RoomID   string              `json:"room_id"`
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

		msg := &types.Message{}
		err = json.Unmarshal(m, msg)
		if err != nil {
			fmt.Printf("error unmarshaling: %v", err)
			continue
		}

		logger.Log.Infof("Received message: %v", msg)

		if err := handleMessage(hub, msg); err != nil {
			errorPayload, _ := json.Marshal(map[string]string{"error": err.Error()})
			hub.SendMessageToClient(&types.Message{Action: "error", Payload: errorPayload, UserID: c.ID})
			fmt.Printf("error handling message: %v", err)
			continue
		}

	}
}

func handleMessage(hub *Hub, msg *types.Message) error {
	action, err := CreateAction(msg, hub)
	if err != nil {
		return fmt.Errorf("error creating action: %v", err)
	}

	err = json.Unmarshal(msg.Payload, action)
	if err != nil {
		return fmt.Errorf("error unmarshalling payload: %v", err)
	}

	return action.Execute()
}

func (c *Client) Disconnect() {
	c.Conn.Close()
	logger.Log.Info("Client disconnected", c.Username)
}
