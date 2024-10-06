package ws

import (
	"encoding/json"
	"fmt"

	"github.com/prodanov17/znk/internal/game"
)

type DealCardsAction struct {
	BaseAction
	Payload json.RawMessage
}

type DealCardsPayload struct {
	Cards []game.Card `json:"cards"`
}

func (a *DealCardsAction) Execute(hub *Hub) error {
	fmt.Println("Received payload:", string(a.Payload))
	if a.RoomID == "" {
		return fmt.Errorf("room id is required")
	}
	game := hub.Room[a.RoomID].Game

	err := game.DealCards()
	if err != nil {
		return fmt.Errorf("failed to deal cards: %w", err)
	}

	for _, team := range game.GameTeam {
		for _, player := range team.Players {
			var dealCardsPayload DealCardsPayload
			fmt.Printf("Player %s hand: %+v\n", player.UserID, player.Hand)
			dealCardsPayload.Cards = player.Hand

			rawPayload, err := json.Marshal(dealCardsPayload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}

			hub.SendMessageToClient(&Message{Action: "deal_cards", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID})

		}

	}

	return nil
}
