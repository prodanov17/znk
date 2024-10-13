package game

type Player struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Hand     []Card `json:"hand"`
}

func NewPlayer(userID string, username string) *Player {
	return &Player{
		UserID:   userID,
		Username: username,
		Hand:     []Card{},
		Score:    0,
	}
}

func (p *Player) AddCard(c *Card) {
	p.Hand = append(p.Hand, *c)
}

func (p *Player) ThrowCard(cardID int) *Card {
	for i, card := range p.Hand {
		if card.ID == cardID {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return &card
		}
	}

	return nil
}

func (p *Player) UpdateScore(score int) {
	p.Score += score
}

func (p *Player) ResetHand() {
	p.Hand = []Card{}
}

func (p *Player) IsEmptyHand() bool {
	return len(p.Hand) == 0
}

func (p *Player) TotalValue() int {
	total := 0
	for _, card := range p.Hand {
		total += card.Value
	}
	return total
}
