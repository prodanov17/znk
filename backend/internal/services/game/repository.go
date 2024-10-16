package game

import (
	"fmt"

	"github.com/prodanov17/znk/internal/services/gamestate"
)

type Repository struct {
	games map[string]*gamestate.Game
}

func NewRepository() *Repository {
	return &Repository{
		games: make(map[string]*gamestate.Game),
	}
}

func (r *Repository) Create(game *gamestate.Game) error {
	r.games[game.ID] = game
	return nil
}
func (r *Repository) FindByID(gameID string) (*gamestate.Game, error) {
	game, ok := r.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}
func (r *Repository) FindAll() ([]*gamestate.Game, error) {
	var games []*gamestate.Game
	for _, game := range r.games {
		games = append(games, game)
	}
	return games, nil
}
func (r *Repository) ClearGame(gameID string) error {
	_, err := r.FindByID(gameID)
	if err != nil {
		return err
	}
	delete(r.games, gameID)
	return nil
}
