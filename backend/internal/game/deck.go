package game

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/prodanov17/znk/internal/config"
)

type Card struct {
	ID    int    `json:"id"`
	Suit  string `json:"suit"`
	Rank  string `json:"rank"`
	Value int    `json:"value"`
	Deck  *Deck  `json:"-"`
}

type Deck struct {
	ID              int `json:"id"`
	Cards           []Card
	CardDefinitions map[int]Card
	Game
}

func NewDeck() *Deck {
	return &Deck{
		ID:              0,
		Cards:           []Card{},
		CardDefinitions: make(map[int]Card),
	}
}

func (d *Deck) Card(id int) Card {
	return d.CardDefinitions[id]
}

func (d *Deck) InitDeck(filepath string) error {
	// Read the JSON file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the JSON data into the deck's Cards slice
	err = json.Unmarshal(data, &d.Cards)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
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
	fmt.Println("Deck after draw", len(d.Cards))
	return &card, nil
}

// compareCards compares two cards and returns true if they have the same rank.
func CompareCards(card1 Card, card2 Card) bool {
	return card1.Rank == card2.Rank
}

func CreateDeck() *Deck {
	deck := NewDeck()
	deck.InitDeck(config.Env.DeckPath)

	for _, card := range deck.Cards {
		deck.CardDefinitions[card.ID] = card
	}

	fmt.Println("Deck created pre-shuffle", len(deck.Cards))
	deck.Shuffle()
	fmt.Println("Deck created", len(deck.Cards))
	return deck
}

func (d *Deck) DreamCard() Card {
	if len(d.Cards) == 0 {
		return Card{}
	}
	return d.Cards[0]
}
