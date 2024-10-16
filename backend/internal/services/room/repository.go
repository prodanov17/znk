package room

import (
	"fmt"

	"github.com/prodanov17/znk/internal/types"
)

type Repository struct {
	rooms map[string]*types.Room
}

func NewRepository() *Repository {
	return &Repository{
		rooms: make(map[string]*types.Room),
	}
}

func (r *Repository) Create(room *types.Room) error {
	r.rooms[room.RoomID] = room
	return nil
}
func (r *Repository) FindByID(roomID string) (*types.Room, error) {
	room, ok := r.rooms[roomID]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}
func (r *Repository) FindAll() ([]*types.Room, error) {
	var rooms []*types.Room
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}
func (r *Repository) AddPlayerToRoom(roomID string, player *types.Player) error {
	room, ok := r.rooms[roomID]
	if !ok {
		return fmt.Errorf("room not found")
	}

	room.Players = append(room.Players, player)

	return nil
}

func (r *Repository) RemovePlayerFromRoom(roomID, playerID string) error {
	room, ok := r.rooms[roomID]
	if !ok {
		return fmt.Errorf("room not found")
	}

	for i, player := range room.Players {
		if player.UserID == playerID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}

	return nil
}

func (r *Repository) ClearRoom(roomID string) error {
	delete(r.rooms, roomID)
	return nil
}
