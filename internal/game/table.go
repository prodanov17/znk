package game

import "fmt"

type Table struct {
	Cards         []Card `json:"cards"`
	IsInitialDeal bool   `json:"initial_deal"`
}

func NewTable() *Table {
	return &Table{
		Cards:         []Card{},
		IsInitialDeal: true,
	}
}

// AddCard adds a card to the table and returns the total value of the cards on the table
// if the card is a match with the first card on the table. If the card is not a match, the
// function returns -1 and the table cards.
func (t *Table) AddCard(card Card) {
	t.Cards = append(t.Cards, card)
}

// ResetTable resets the table to an empty state.
func (t *Table) ResetTable() {
	t.Cards = []Card{}
	t.IsInitialDeal = true
}

// TotalValue returns the total value of the cards on the table.
func (t *Table) TotalValue() int {
	total := 0
	for _, card := range t.Cards {
		total += card.Value
	}
	return total
}

func (t *Table) InitialDeal(deck *Deck) {
	for i := 0; i < 4; i++ {
		card, _ := deck.DrawCard()
		t.Cards = append(t.Cards, *card)
	}
}

func (t *Table) TopCard() (Card, error) {
	if len(t.Cards) == 0 {
		return Card{}, fmt.Errorf("no cards on the table")
	}
	return t.Cards[len(t.Cards)-1], nil
}
