package ws

import (
	"fmt"
)

type Action interface {
	Execute() error
}

type BaseAction struct {
	Hub    *Hub
	UserID string
	RoomID string
}

func (a *BaseAction) Execute() error {
	return nil
}

func NewBaseAction(hub *Hub, userID, roomID string) BaseAction {
	return BaseAction{
		Hub:    hub,
		UserID: userID,
		RoomID: roomID,
	}
}

var Actions = map[string]func(message *Message, hub *Hub) Action{
	"message": func(message *Message, hub *Hub) Action {
		return &MessageAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"deal_cards": func(message *Message, hub *Hub) Action {
		return &DealCardsAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"throw_card": func(message *Message, hub *Hub) Action {
		return &ThrowCardAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"start_game": func(message *Message, hub *Hub) Action {
		return &StartGameAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"change_team": func(message *Message, hub *Hub) Action {
		return &ChangeTeamAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"get_teams": func(message *Message, hub *Hub) Action {
		return &GetTeamsAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"game_state": func(message *Message, hub *Hub) Action {
		return &GetGameStateAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
	"get_dealer": func(message *Message, hub *Hub) Action {
		return &GetDealerAction{
			BaseAction: NewBaseAction(hub, message.UserID, message.RoomID),
			Payload:    message.Payload,
		}
	},
}

func CreateAction(message *Message, hub *Hub) (Action, error) {
	constructor, ok := Actions[message.Action]
	if !ok {
		return nil, fmt.Errorf("action type %s not found", message.Action)
	}
	return constructor(message, hub), nil
}
