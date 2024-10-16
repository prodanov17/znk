package game

import (
	"github.com/prodanov17/znk/internal/services/gamestate"
	"github.com/prodanov17/znk/internal/types"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGame(gamePayload *types.CreateGamePayload) (*gamestate.Game, error) {
	game := gamestate.NewGame(gamePayload.RoomID, gamePayload.UserID)
	err := s.repo.Create(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}
func (s *Service) GetGameByID(id string) (*gamestate.Game, error) {
	return s.repo.FindByID(id)
}
func (s *Service) GetGames() ([]*gamestate.Game, error) {
	return s.repo.FindAll()
}
func (s *Service) AddUserToGame(gameID, userID, username string) error {
	game, err := s.repo.FindByID(gameID)
	if err != nil {
		return err
	}
	player := &gamestate.Player{UserID: userID, Username: username}
	game.AddPlayer(player)
	return nil
}
func (s *Service) RemoveUserFromGame(gameID string, userID string) error {
	game, err := s.repo.FindByID(gameID)
	if err != nil {
		return err
	}
	game.RemovePlayer(userID)
	return nil
}
