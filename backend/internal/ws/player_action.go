package ws

import (
	"encoding/json"
	"fmt"

	"github.com/prodanov17/znk/internal/services/gamestate"
	"github.com/prodanov17/znk/internal/types"
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
	CardID           int              `json:"card_id"`
	Card             gamestate.Card   `json:"card"`
	TakeCards        bool             `json:"take_cards"`
	Value            int              `json:"value"`
	UserID           string           `json:"user_id"`
	Username         string           `json:"username"`
	RoomID           string           `json:"room_id"`
	TableCards       []gamestate.Card `json:"table_cards"`
	TableValue       int              `json:"table_value"`
	PlayersCardCount map[string]int   `json:"players_card_count"`
	PlayerHand       []gamestate.Card `json:"player_hand"`
	Playing          bool             `json:"playing"`
	NextTurn         string           `json:"next_turn_id"`
}

func (a *ThrowCardAction) Execute() error {
	var throwCardPayload ThrowCardPayload
	err := json.Unmarshal(a.Payload, &throwCardPayload)
	if err != nil {
		return err
	}

	game, err := a.Hub.RoomService().GameService().GetGameByID(a.RoomID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	value, err := game.PlayCard(a.UserID, throwCardPayload.CardID)
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
			Username:         player.Username,
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

		message := &types.Message{Action: "card_played", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID}
		a.Hub.SendMessageToClient(message)
	}

	if game.RoundOver() {
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

			message := &types.Message{Action: "game_ended", Payload: rawPayload, RoomID: a.RoomID}
			a.Hub.BroadcastMessage(message)
			a.Hub.SendMessageToClient(message)
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
			message := &types.Message{Action: "round_over", Payload: rawPayload, RoomID: a.RoomID, UserID: "0"}
			a.Hub.BroadcastMessage(message)

		}
	}

	return nil
}
