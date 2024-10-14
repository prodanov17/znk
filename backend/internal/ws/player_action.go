package ws

import (
	"encoding/json"
	"fmt"

	"github.com/prodanov17/znk/internal/game"
)

type ThrowCardAction struct {
	BaseAction
	Payload json.RawMessage
}

type ThrowCardPayload struct {
	CardID int `json:"card_id"`
}

type PlayerCardCount struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}
type ThrowCardResponse struct {
	CardID           int            `json:"card_id"`
	Card             game.Card      `json:"card"`
	TakeCards        bool           `json:"take_cards"`
	Value            int            `json:"value"`
	UserID           string         `json:"user_id"`
	Username         string         `json:"username"`
	RoomID           string         `json:"room_id"`
	TableCards       []game.Card    `json:"table_cards"`
	TableValue       int            `json:"table_value"`
	PlayersCardCount map[string]int `json:"players_card_count"`
	PlayerHand       []game.Card    `json:"player_hand"`
	Playing          bool           `json:"playing"`
	NextTurn         string         `json:"next_turn_id"`
}

func (a *ThrowCardAction) Execute(hub *Hub) error {
	var throwCardPayload ThrowCardPayload
	err := json.Unmarshal(a.Payload, &throwCardPayload)
	if err != nil {
		return err
	}

	game := hub.Room[a.RoomID].Game

	value, err := game.PlayCard(a.UserID, throwCardPayload.CardID)
	fmt.Println("value", value)
	if err != nil {
		return err
	}
	game.GameState.AdvanceTurn()
	nextTurn, err := game.GameState.NextTurn(game.GameTeam)

	for _, player := range game.Players() {
		throwCardResponse := ThrowCardResponse{
			CardID:           throwCardPayload.CardID,
			Card:             game.Deck.Card(throwCardPayload.CardID),
			TakeCards:        value != -1,
			Value:            value,
			UserID:           a.UserID,
			Username:         hub.Clients[a.UserID].Username,
			RoomID:           a.RoomID,
			TableCards:       game.TableCards(),
			TableValue:       game.Table.TotalValue(),
			PlayersCardCount: game.PlayersCardCount(),
			PlayerHand:       game.PlayerHand(player.UserID),
			Playing:          !game.CanDealCards(),
			NextTurn:         nextTurn.UserID,
		}

		rawPayload, err := json.Marshal(throwCardResponse)
		if err != nil {
			return err
		}

		message := &Message{Action: "card_played", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID}
		hub.SendMessageToClient(message)
	}

	if game.RoundOver() {
		fmt.Println("Round over")
		winner, err := game.EndRound()
		if err != nil {
			return err
		}

		if winner != nil { //check this
			roundOverPayload := map[string]interface{}{
				"winner_team": winner.ID,
				"teams":       game.GameTeam,
			}

			rawPayload, err := json.Marshal(roundOverPayload)
			if err != nil {
				return err
			}

			message := &Message{Action: "game_ended", Payload: rawPayload, RoomID: a.RoomID}
			hub.BroadcastMessage(message)
			hub.SendMessageToClient(message)
			return nil
		} else {

			dealer, err := game.GameState.Dealer(game.GameTeam)
			if err != nil {
				return fmt.Errorf("failed to get dealer: %w", err)
			}

			roundOverPayload := map[string]interface{}{"dealer_id": dealer.UserID, "teams": game.GameTeam, "last_capture_id": game.LastCapture().ID}
			rawPayload, err := json.Marshal(roundOverPayload)
			if err != nil {
				return err
			}
			message := &Message{Action: "round_over", Payload: rawPayload, RoomID: a.RoomID, UserID: "0"}
			hub.BroadcastMessage(message)

		}
	}

	fmt.Println("Next turn:", nextTurn.UserID)

	return nil
}
