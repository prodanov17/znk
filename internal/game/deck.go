package game

import (
	"fmt"
	"math/rand"
)

type Card struct {
	ID    int    `json:"id"`
	Suit  string `json:"suit"`
	Rank  string `json:"rank"`
	Value int    `json:"value"`
	Deck  *Deck  `json:"-"`
}

type Deck struct {
	ID    int `json:"id"`
	Cards []Card
	Game
}

func NewDeck() *Deck {
	return &Deck{
		ID:    0,
		Cards: []Card{},
	}
}

func (d *Deck) InitDeck() {
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for _, suit := range suits {
		for _, rank := range ranks {
			card := Card{
				ID:    len(d.Cards),
				Suit:  suit,
				Rank:  rank,
				Value: 0,
			}
			d.Cards = append(d.Cards, card)
		}
	}
}

func (d *Deck) Shuffle() {
	for i := range d.Cards {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

func (d *Deck) DrawCard() (*Card, error) {
	if len(d.Cards) == 0 {
		return nil, fmt.Errorf("No more cards in deck")
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return &card, nil
}
