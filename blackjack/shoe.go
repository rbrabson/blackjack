package blackjack

import "github.com/rbrabson/blackjack/deck"

// Shoe represents a collection of decks of cards used in blackjack.
type Shoe struct {
	reshufflePercentage int
	numDecks            int
	deck.Cards
}

// NewShoe creates a new Shoe with the specified number of decks and reshuffle percentage.
func NewShoe(numDecks int, reshufflePercentage int) Shoe {
	shoe := Shoe{
		reshufflePercentage: reshufflePercentage,
		numDecks:            numDecks,
	}
	shoe.ReShuffle()

	return shoe
}

// StartNewRound checks if the shoe needs to be reshuffled based on the reshuffle percentage.
// If it does, then all cards are added back to the shoe and shuffled again.
func (s *Shoe) StartNewRound() {
	cut := (s.reshufflePercentage * len(s.Cards.Cards)) / 100
	if len(s.Cards.Cards) <= cut {
		s.ReShuffle()
	}
}

// ReShuffle adds all cards from the specified number of decks to the shoe and shuffles them.
func (s *Shoe) ReShuffle() {
	s.Cards = deck.Cards{Cards: []deck.Card{}}
	for i := 0; i < s.numDecks; i++ {
		deck := deck.New()
		s.Cards.Cards = append(s.Cards.Cards, deck.Cards.Cards...)
	}
	s.Shuffle()
}
