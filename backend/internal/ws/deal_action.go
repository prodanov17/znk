package ws

import (
	"encoding/json"
	"fmt"

	"github.com/prodanov17/znk/internal/services/gamestate"
	"github.com/prodanov17/znk/internal/types"
)

type DealCardsAction struct {
	BaseAction
	Payload json.RawMessage
}

type DealCardsPayload struct {
	Cards            []gamestate.Card `json:"cards"`
	TableCards       []gamestate.Card `json:"table_cards"`
	TableValue       int              `json:"table_value"`
	PlayersCardCount map[string]int   `json:"players_card_count"`
	NextTurnId       string           `json:"next_turn_id"`
	DreamCard        gamestate.Card   `json:"dream_card"`
}

func (a *DealCardsAction) Execute() error {
	if a.RoomID == "" {
		return fmt.Errorf("room id is required")
	}
	game, err := a.Hub.RoomService().GameService().GetGameByID(a.RoomID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

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

			a.Hub.SendMessageToClient(&types.Message{Action: "deal_cards", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID})

		}

	}

	return nil
}
