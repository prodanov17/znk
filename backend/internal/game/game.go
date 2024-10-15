package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/prodanov17/znk/pkg/logger"
)

type Game struct {
	sync.Mutex
	ID        string         `json:"id"`
	OwnerID   string         `json:"owner_id"`
	CreatedAt string         `json:"created_at"`
	State     string         `json:"state"`
	GameTeam  []GameTeam     `json:"-"`
	Deck      *Deck          `json:"-"`
	GameState *GameState     `json:"-"`
	Rules     map[string]int `json:"-"`
	Table     *Table         `json:"-"`
}

func NewGame(id, ownerID string) *Game {
	deck := CreateDeck()
	rules := DefaultRules()

	teamA := NewGameTeam(1, id)
	teamB := NewGameTeam(2, id)

	GameState := NewGameState()
	Table := NewTable()

	return &Game{
		ID:        id,
		OwnerID:   ownerID,
		GameTeam:  []GameTeam{*teamA, *teamB},
		Deck:      deck,
		Rules:     rules,
		State:     "waiting",
		GameState: GameState,
		Table:     Table,
		CreatedAt: time.Now().String(),
	}
}

func (g *Game) ResetGame() {
	g.GameTeam[0].Score = 0
	g.GameTeam[1].Score = 0
	g.GameState.Reset()
	g.Table.ResetTable()
	g.Deck = CreateDeck()
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
	for i := range g.GameTeam {
		for j := range g.GameTeam[i].Players {
			if g.GameTeam[i].Players[j].UserID == playerID {
				g.GameTeam[i].RemovePlayer(g.GameTeam[i].Players[j])
				return nil
			}
		}
	}

	return fmt.Errorf("Player not found") //rejoining with the same id creates a bug
}

func (g *Game) StartGame() error {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if len(g.GameTeam[0].Players) < 2 || len(g.GameTeam[1].Players) < 2 {
		return fmt.Errorf("Not enough players in teams, have %d need 4", len(g.GameTeam[0].Players)+len(g.GameTeam[1].Players))
	}

	if g.State == "started" {
		return fmt.Errorf("Game already started")
	}
	g.State = "started"

	g.ResetGame()

	return nil
}

func (g *Game) EndGame() {
	g.State = "ended"
}

func (g *Game) DealCards() error {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	fmt.Println("Dealing cards", g.State)
	if g.State != "started" {
		return fmt.Errorf("Game is not started")
	}
	fmt.Println("Can deal cards", g.CanDealCards())
	if !g.CanDealCards() {
		return fmt.Errorf("Cannot deal cards")
	}

	if g.Table.IsInitialDeal {
		g.Table.ResetTable()
		g.Table.InitialDeal(g.Deck)
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("Dealing round %d\n", i+1)
		for _, team := range g.GameTeam {
			for _, player := range team.Players {
				fmt.Printf("Drawing card for player %s\n", player.Username)
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

// CanDealCards checks if all players have an empty hand
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

// Player returns a player by their ID
func (g *Game) Player(playerID string) (*Player, error) {
	for _, team := range g.GameTeam {
		for i := range team.Players {
			if team.Players[i].UserID == playerID {
				return team.Players[i], nil // Return pointer to the player
			}
		}
	}

	return nil, fmt.Errorf("Player not found")
}

// PlayCard plays a card from a player's hand, adds it to the table and checks if a capture is possible
// If a capture is possible, the cards are captured and the player's and team's score is updated
// The function returns the value of the capture if it happened
func (g *Game) PlayCard(playerID string, cardID int) (int, error) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if !g.IsPlayerTurn(playerID) {
		return -1, fmt.Errorf("Not your turn")
	}
	player, err := g.Player(playerID)
	if err != nil {
		return -1, err
	}

	card := player.ThrowCard(cardID)
	if card == nil {
		return -1, fmt.Errorf("Card not found")
	}

	g.Table.AddCard(*card)

	captureValue := g.checkCapture()

	if captureValue != -1 {
		team := g.PlayerTeam(playerID)
		team.CaptureCards(g.Table.Cards)

		team.UpdateScore(captureValue)
		player.UpdateScore(captureValue)

		g.Table.ResetTable()
	}

	g.Table.IsInitialDeal = false

	return captureValue, nil
}

func (g *Game) RoundOver() bool {
	if len(g.Deck.Cards) != 0 {
		return false
	}

	for _, team := range g.GameTeam {
		for _, player := range team.Players {
			fmt.Printf("Player hand length %d for player %s", len(player.Hand), player.UserID)
			if len(player.Hand) != 0 {
				return false
			}
		}
	}

	return true
}

func (g *Game) EndRound() (*GameTeam, error) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if len(g.Deck.Cards) != 0 {
		return nil, fmt.Errorf("Cannot end round, deck is not empty")
	}
	fmt.Println("Ending round")
	g.GameState.AdvanceDealer()
	g.GameState.TurnIdx = g.GameState.DealerIdx + 1

	tableValue := g.Table.TotalValue()
	g.LastCapture().UpdateScore(tableValue)
	g.LastCapture().CaptureCards(g.Table.Cards)

	if len(g.GameTeam[0].CapturedCards) > len(g.GameTeam[1].CapturedCards) {
		g.GameTeam[0].UpdateScore(4)
	} else if len(g.GameTeam[1].CapturedCards) > len(g.GameTeam[0].CapturedCards) {
		g.GameTeam[1].UpdateScore(4)
	} else {
		g.GameTeam[0].UpdateScore(2)
		g.GameTeam[1].UpdateScore(2)
	}

	g.Table.ResetTable()

	winner := g.Winner()
	if winner != nil {
		g.EndGame()
		return winner, nil
	}

	g.Deck = CreateDeck()
	g.GameTeam[0].ResetCapturedCards()
	g.GameTeam[1].ResetCapturedCards()

	return nil, nil
}

func (g *Game) Winner() *GameTeam {
	winningScore := g.Rules["winning_score"]

	team1Score := g.GameTeam[0].Score
	team2Score := g.GameTeam[1].Score

	fmt.Println("Team 1 score", team1Score)
	fmt.Println("Team 2 score", team2Score)

	if team1Score > winningScore || team2Score > winningScore {
		if team1Score > team2Score {
			return &g.GameTeam[0]
		} else if team2Score > team1Score {
			return &g.GameTeam[1]
		}
	}

	return nil
}

func (g *Game) IsPlayerTurn(playerID string) bool {
	turn, err := g.GameState.NextTurn(g.GameTeam)

	if err != nil {
		logger.Log.Info("Error getting next turn:", err)
		return false
	}

	return turn.UserID == playerID
}

func (g *Game) TableCards() []Card {
	return g.Table.Cards
}

func (g *Game) PlayerHand(playerID string) []Card {
	player, err := g.Player(playerID)
	if err != nil {
		return []Card{}
	}

	return player.Hand
}

func (g *Game) Players() []*Player {
	players := []*Player{}
	for _, team := range g.GameTeam {
		players = append(players, team.Players...)
	}

	return players
}

func (g *Game) PlayerTeam(playerID string) *GameTeam {
	for i := range g.GameTeam {
		for _, player := range g.GameTeam[i].Players {
			if player.UserID == playerID {
				return &g.GameTeam[i] // Return the actual team pointer
			}
		}
	}
	return nil
}

func (g *Game) ChangeTeam(playerID string) error {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if g.State == "started" {
		return fmt.Errorf("cannot change team at this state")
	}
	player, err := g.Player(playerID)
	if err != nil {
		return err
	}

	team := g.PlayerTeam(playerID)

	if team.ID == 1 {
		g.GameTeam[0].RemovePlayer(player)
		g.GameTeam[1].AddPlayer(player)
	} else {
		g.GameTeam[1].RemovePlayer(player)
		g.GameTeam[0].AddPlayer(player)
	}

	return nil
}

func (g *Game) LastCapture() *GameTeam {
	dealer, err := g.GameState.Dealer(g.GameTeam)
	if err != nil {
		return &g.GameTeam[0]
	}

	dealerTeam := g.PlayerTeam(dealer.UserID)

	if len(g.Table.Cards)%2 == 0 {
		return dealerTeam
	}

	for i := range g.GameTeam {
		if &g.GameTeam[i] != dealerTeam {
			return &g.GameTeam[i]
		}
	}

	return &g.GameTeam[0]
}

func (g *Game) checkCapture() int {
	topCard, err := g.Table.TopCard()
	if err != nil {
		return -1
	}

	if len(g.Table.Cards) == 2 {
		if topCard.Rank == g.Table.Cards[0].Rank {
			return g.Table.TotalValue() + 10
		}
	}

	if topCard.Rank == "J" && len(g.Table.Cards) > 1 {
		return g.Table.TotalValue()
	}

	if len(g.Table.Cards) < 2 {
		return -1
	}
	if g.Table.IsInitialDeal {
		if topCard.Rank == g.Table.Cards[0].Rank {
			return g.Table.TotalValue()
		}
	} else {
		if topCard.Rank == g.Table.Cards[len(g.Table.Cards)-2].Rank {
			return g.Table.TotalValue()
		}
	}

	return -1
}

func (g *Game) PlayersCardCount() map[string]int {
	counts := map[string]int{}
	for _, team := range g.GameTeam {
		for _, player := range team.Players {
			counts[player.UserID] = len(player.Hand)
		}
	}

	return counts
}

func (g *Game) UpdateRule(rule string, value int) {
	g.Rules[rule] = value
}
