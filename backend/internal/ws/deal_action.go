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
	Cards            []game.Card    `json:"cards"`
	TableCards       []game.Card    `json:"table_cards"`
	TableValue       int            `json:"table_value"`
	PlayersCardCount map[string]int `json:"players_card_count"`
	NextTurnId       string         `json:"next_turn_id"`
	DreamCard        game.Card      `json:"dream_card"`
}

func (a *DealCardsAction) Execute(hub *Hub) error {
	if a.RoomID == "" {
		return fmt.Errorf("room id is required")
	}
	game := hub.Room[a.RoomID].Game

	dealer, err := game.GameState.Dealer(game.GameTeam)
	if err != nil {
		return fmt.Errorf("failed to get dealer: %w", err)
	}

	if dealer.UserID != a.UserID {
		return fmt.Errorf("only the dealer can deal cards")
	}

	err = game.DealCards()
	if err != nil {
		return fmt.Errorf("failed to deal cards: %w", err)
	}

	nextTurn, err := game.GameState.NextTurn(game.GameTeam)

	for _, team := range game.GameTeam {
		for _, player := range team.Players {
			var dealCardsPayload DealCardsPayload
			fmt.Printf("Player %s hand: %+v\n", player.UserID, player.Hand)
			dealCardsPayload.Cards = player.Hand
			dealCardsPayload.TableCards = game.TableCards()
			dealCardsPayload.TableValue = game.Table.TotalValue()
			dealCardsPayload.NextTurnId = nextTurn.UserID
			dealCardsPayload.PlayersCardCount = game.PlayersCardCount()
			dealCardsPayload.DreamCard = game.Deck.DreamCard()

			rawPayload, err := json.Marshal(dealCardsPayload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}

			hub.SendMessageToClient(&Message{Action: "deal_cards", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID})

		}

	}

	return nil
}
