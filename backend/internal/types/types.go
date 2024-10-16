package types

import (
	"encoding/json"
	"time"

	"github.com/prodanov17/znk/internal/services/gamestate"
)

type UserService interface {
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	LoginUser(userPayload *LoginUserPayload) (string, error)
	RegisterUser(userPayload *RegisterUserPayload) (string, error)
	UpdateUser(id int, userPayload *UpdateUserPayload) (*User, error)
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
}

type Player struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
}

type Room struct {
	RoomID    string    `json:"room_id"`
	Players   []*Player `json:"clients"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

type RoomService interface {
	GameService() GameService
	CreateRoom(roomPayload *CreateRoomPayload) (*Room, error)
	GetRoomByID(id string) (*Room, error)
	GetRooms() ([]*Room, error)
	AddUserToRoom(roomID, userID, username string) error
	RemoveUserFromRoom(roomID, userID string) error
	ClearRoom(roomID string) error
	GetPlayerById(roomID, playerID string) (*Player, error)
}

type RoomRepository interface {
	Create(room *Room) error
	FindByID(roomID string) (*Room, error)
	FindAll() ([]*Room, error)
	AddPlayerToRoom(roomID string, player *Player) error
	RemovePlayerFromRoom(roomID string, playerID string) error
	ClearRoom(roomID string) error
}

type GameService interface {
	CreateGame(gamePayload *CreateGamePayload) (*gamestate.Game, error)
	GetGameByID(id string) (*gamestate.Game, error)
	GetGames() ([]*gamestate.Game, error)
	AddUserToGame(gameID, userID, username string) error
	RemoveUserFromGame(gameID string, userID string) error
}

type GameRepository interface {
	Create(game *gamestate.Game) error
	FindByID(gameID string) (*gamestate.Game, error)
	FindAll() ([]*gamestate.Game, error)
	AddPlayerToGame(gameID string, player *Player) error
	RemovePlayerFromGame(gameID string, playerID string) error
	ClearGame(gameID string) error
}

type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
	UserID  string          `json:"user_id"`
	RoomID  string          `json:"room_id"`
}

// type HubInterface interface {
// 	BroadcastMessage(message *Message)
// 	SendMessageToClient(message *Message)
// 	UnregisterClient(clientID string)
// 	RegisterClient(client *Client)
// 	RoomService() RoomService
// 	Clients() map[string]*Client
// }
