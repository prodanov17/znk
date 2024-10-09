package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/prodanov17/znk/internal/game"
)

type StartGameAction struct {
	BaseAction
	Payload json.RawMessage
}

type StartGamePayload struct {
	Rules map[string]string `json:"rules"`
}

func (a *StartGameAction) Execute(hub *Hub) error {
	game := hub.Room[a.RoomID].Game

	fmt.Println("gamerule", game.Rules)
	err := game.StartGame()
	if err != nil {
		return err
	}

	var startGamePayload StartGamePayload

	err = json.Unmarshal(a.Payload, &startGamePayload)
	if err != nil {
		log.Println("Failed to unmarshal payload. Using default rules")
	}

	for key, value := range startGamePayload.Rules {
		fmt.Println("key", key, "value", value)
		intValue, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			continue
		}
		game.UpdateRule(key, intValue)
	}

	fmt.Println("gamerule", game.Rules)
	dealer, err := game.GameState.Dealer(game.GameTeam)
	if err != nil {
		return err
	}
	tableCards := map[string]interface{}{"dealer_id": dealer.UserID}

	rawPayload, err := json.Marshal(tableCards)
	if err != nil {
		return err
	}

	message := &Message{Action: "game_started", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	hub.BroadcastMessage(message)
	hub.SendMessageToClient(message)

	return nil
}

type ChangeTeamAction struct {
	BaseAction
	Payload json.RawMessage
}

func (a *ChangeTeamAction) Execute(hub *Hub) error {
	game := hub.Room[a.RoomID].Game

	err := game.ChangeTeam(a.UserID)
	if err != nil {
		return err
	}

	playerTeam := game.PlayerTeam(a.UserID)
	if playerTeam == nil {
		return nil
	}
	var playerPayload = map[string]interface{}{"id": a.UserID, "team_id": playerTeam.ID, "teams": game.GameTeam}
	rawPayload, _ := json.Marshal(playerPayload)

	message := &Message{Action: "team_changed", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	hub.BroadcastMessage(message)
	hub.SendMessageToClient(message)

	return nil
}

type GetTeamsAction struct {
	BaseAction
	Payload json.RawMessage
}

func (a *GetTeamsAction) Execute(hub *Hub) error {
	game := hub.Room[a.RoomID].Game

	teams := game.GameTeam
	var teamPayload = map[string]interface{}{"teams": teams}
	rawPayload, _ := json.Marshal(teamPayload)

	message := &Message{Action: "teams", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	hub.SendMessageToClient(message)

	return nil
}

type GetDealerAction struct {
	BaseAction
	Payload json.RawMessage
}

func (a *GetDealerAction) Execute(hub *Hub) error {
	game := hub.Room[a.RoomID].Game

	dealer, err := game.GameState.Dealer(game.GameTeam)
	if err != nil {
		return err
	}

	var dealerPayload = map[string]interface{}{"dealer_id": dealer.UserID}
	rawPayload, _ := json.Marshal(dealerPayload)

	message := &Message{Action: "dealer", Payload: rawPayload, RoomID: a.RoomID, UserID: a.UserID}
	hub.SendMessageToClient(message)

	return nil
}

type GetGameStateAction struct {
	BaseAction
	Payload json.RawMessage
}

type GameStatePayload struct {
	State            string          `json:"state"`
	Teams            []game.GameTeam `json:"teams"`
	PlayerHand       []game.Card     `json:"player_hand"`
	TableCards       []game.Card     `json:"table_cards"`
	TableValue       int             `json:"table_value"`
	PlayersCardCount map[string]int  `json:"players_card_count"`
	NextTurnId       string          `json:"next_turn_id"`
	DreamCard        game.Card       `json:"dream_card"`
	DealerID         string          `json:"dealer_id"`
	Playing          bool            `json:"playing"`
	GameInfo         map[string]int  `json:"game_info"`
}

func (a *GetGameStateAction) Execute(hub *Hub) error {
	game := hub.Room[a.RoomID].Game

	nextTurn, _ := game.GameState.NextTurn(game.GameTeam)
	dealer, _ := game.GameState.Dealer(game.GameTeam)

	for _, team := range game.GameTeam {
		for _, player := range team.Players {
			gameStatePayload := GameStatePayload{
				State:            game.State,
				Teams:            game.GameTeam,
				TableCards:       game.TableCards(),
				PlayerHand:       game.PlayerHand(player.UserID),
				TableValue:       game.Table.TotalValue(),
				NextTurnId:       nextTurn.UserID,
				PlayersCardCount: game.PlayersCardCount(),
				DreamCard:        game.Deck.DreamCard(),
				DealerID:         dealer.UserID,
				Playing:          game.CanDealCards(),
				GameInfo:         game.Rules,
			}

			rawPayload, err := json.Marshal(gameStatePayload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}

			hub.SendMessageToClient(&Message{Action: "game_state", Payload: rawPayload, RoomID: a.RoomID, UserID: player.UserID})

		}

	}

	return nil
}
