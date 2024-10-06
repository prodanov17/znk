package game

import (
	"fmt"
	"time"
)

type Game struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
	GameTeam  []GameTeam
	Deck      *Deck
	Rules     []GameRule
}

func NewGame(id, ownerID string) *Game {
	deck := NewDeck()
	deck.InitDeck()
	deck.Shuffle()

	rules := DefaultRules()

	teamA := NewGameTeam(1, id)
	teamB := NewGameTeam(2, id)

	return &Game{
		ID:        id,
		OwnerID:   ownerID,
		GameTeam:  []GameTeam{*teamA, *teamB},
		Deck:      deck,
		Rules:     rules,
		State:     "waiting",
		CreatedAt: time.Now().String(),
	}
}

func (g *Game) AddPlayer(p *Player) error {
	if len(g.GameTeam[0].Players) < 2 {
		g.GameTeam[0].AddPlayer(p)
	} else if len(g.GameTeam[1].Players) < 2 {
		g.GameTeam[1].AddPlayer(p)
	} else {
		return fmt.Errorf("Game is full")
	}

	return nil
}

func (g *Game) RemovePlayer(playerID string) error {
	for _, team := range g.GameTeam {
		for i, player := range team.Players {
			if player.UserID == playerID {
				team.Players = append(team.Players[:i], team.Players[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("Player not found")
}

func (g *Game) StartGame() error {
	if len(g.GameTeam[0].Players) < 2 || len(g.GameTeam[1].Players) < 2 {
		return fmt.Errorf("Not enough players")
	}
	g.State = "started"

	return nil
}

func (g *Game) EndGame() {
	g.State = "ended"
}

func (g *Game) DealCards() error {
	if g.State != "started" {
		return fmt.Errorf("Game is not started")
	}
	if !g.CanDealCards() {
		return fmt.Errorf("Cannot deal cards")
	}

	for i := 0; i < 4; i++ {
		for _, team := range g.GameTeam {
			for _, player := range team.Players {
				fmt.Println("Player:", player.Hand)
				card, err := g.Deck.DrawCard()
				if err != nil {
					fmt.Println("Error drawing card:", err)
					return err
				}

				player.AddCard(card)
			}
		}
	}

	return nil
}

func (g *Game) CanDealCards() bool {
	for _, team := range g.GameTeam {
		for _, player := range team.Players {
			if !player.IsEmptyHand() {
				return false
			}
		}
	}

	return true
}
