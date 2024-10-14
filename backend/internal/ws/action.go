package ws

import (
	"encoding/json"
	"fmt"
)

type Action interface {
	Execute(hub *Hub) error
}

type BaseAction struct {
	UserID string
	RoomID string
}

func (a *BaseAction) Execute(hub *Hub) error {
	return nil
}

func NewBaseAction(userID, roomID string) BaseAction {
	return BaseAction{
		UserID: userID,
		RoomID: roomID,
	}
}

var Actions = map[string]func(userID, roomID string, payload json.RawMessage) Action{
	"example": func(userID, roomID string, payload json.RawMessage) Action {
		return &ExampleAction{BaseAction: BaseAction{
			UserID: userID,
			RoomID: roomID,
		}}
	},
	"message": func(userID, roomID string, payload json.RawMessage) Action {
		return &MessageAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"deal_cards": func(userID, roomID string, payload json.RawMessage) Action {
		return &DealCardsAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"throw_card": func(userID, roomID string, payload json.RawMessage) Action {
		return &ThrowCardAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"start_game": func(userID, roomID string, payload json.RawMessage) Action {
		return &StartGameAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"change_team": func(userID, roomID string, payload json.RawMessage) Action {
		return &ChangeTeamAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"get_teams": func(userID, roomID string, payload json.RawMessage) Action {
		return &GetTeamsAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"game_state": func(userID, roomID string, payload json.RawMessage) Action {
		return &GetGameStateAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
	"get_dealer": func(userID, roomID string, payload json.RawMessage) Action {
		return &GetDealerAction{
			BaseAction: NewBaseAction(userID, roomID),
			Payload:    payload,
		}
	},
}

func CreateAction(actionType, gameID, userID string, payload json.RawMessage) (Action, error) {
	constructor, ok := Actions[actionType]
	if !ok {
		return nil, fmt.Errorf("action type %s not found", actionType)
	}
	return constructor(userID, gameID, payload), nil
}
