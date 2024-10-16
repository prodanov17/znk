package room

import (
	"time"

	"github.com/prodanov17/znk/internal/types"
)

type Service struct {
	gameService types.GameService
	repo        types.RoomRepository
}

func NewService(gameService types.GameService, repo types.RoomRepository) *Service {
	return &Service{
		gameService: gameService,
		repo:        repo,
	}
}

func (s *Service) GameService() types.GameService {
	return s.gameService
}

// CreateRoom creates a new room
func (s *Service) CreateRoom(roomPayload *types.CreateRoomPayload) (*types.Room, error) {
	room := &types.Room{
		RoomID:    roomPayload.RoomID,
		CreatedBy: roomPayload.UserID,
		CreatedAt: time.Now(),
		Players:   []*types.Player{},
	}

	if err := s.repo.Create(room); err != nil {
		return nil, err
	}

	gamePayload := &types.CreateGamePayload{
		RoomID: roomPayload.RoomID,
		UserID: roomPayload.UserID,
	}
	s.gameService.CreateGame(gamePayload)

	return room, nil
}

func (s *Service) GetRoomByID(roomID string) (*types.Room, error) {
	return s.repo.FindByID(roomID)
}

func (s *Service) GetRooms() ([]*types.Room, error) {
	return s.repo.FindAll()
}

func (s *Service) AddUserToRoom(roomID, userID, username string) error {
	// add to gameservice
	player := &types.Player{
		UserID:   userID,
		Username: username,
	}

	err := s.gameService.AddUserToGame(roomID, userID, username)
	if err != nil {
		return err
	}
	return s.repo.AddPlayerToRoom(roomID, player)
}

func (s *Service) RemoveUserFromRoom(roomID string, playerID string) error {
	// remove from gameservice
	s.gameService.RemoveUserFromGame(roomID, playerID)
	return s.repo.RemovePlayerFromRoom(roomID, playerID)
}

func (s *Service) ClearRoom(roomID string) error {
	game, err := s.gameService.GetGameByID(roomID)
	if err != nil {
		return err
	}
	game.ClearGame()
	return s.repo.ClearRoom(roomID)
}

func (s *Service) GetPlayerById(roomID, playerID string) (*types.Player, error) {
	room, err := s.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	for _, player := range room.Players {
		if player.UserID == playerID {
			return player, nil
		}
	}
	return nil, nil
}
