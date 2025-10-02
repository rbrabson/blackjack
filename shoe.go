package blackjack

import (
	"fmt"

	"github.com/rbrabson/cards"
)

const (
	CutCardPenetration = 0.75 // CutCardPenetration is the fraction of the shoe dealt before reshuffling
)

// Shoe wraps the cards.Shoe with blackjack-specific functionality
type Shoe struct {
	cards    cards.Shoe // shoe is the set of cards to be dealt
	numDecks int        // numDecdks is the number of decks in the shoe
	cutCard  int        // Position where cut card is placed (reshuffle point)
}

// NewShoe creates a new blackjack shoe with the specified number of decks
func NewShoe(numDecks int) *Shoe {
	s := &Shoe{
		numDecks: max(1, numDecks),
	}
	s.Reshuffle()

	return s
}

// Draw deals a card from the shoe
func (s *Shoe) Draw() (cards.Card, error) {
	if s.IsEmpty() {
		return cards.Card{}, fmt.Errorf("shoe is empty")
	}

	return s.cards.Draw(), nil
}

// IsEmpty returns true if the shoe is empty
func (s *Shoe) IsEmpty() bool {
	return len(s.cards) == 0
}

// NeedsReshuffle returns true if the cut card has been reached
func (s *Shoe) NeedsReshuffle() bool {
	return len(s.cards) <= (s.numDecks*52 - s.cutCard)
}

// CardsRemaining returns the number of cards left in the shoe
func (s *Shoe) CardsRemaining() int {
	return len(s.cards)
}

// Reshuffle creates a new shuffled shoe with the same number of decks
func (s *Shoe) Reshuffle() {
	s.cards = cards.NewShoe(s.numDecks)
	s.cards.Shuffle()

	// Reset cut card position
	s.cutCard = int(float64(len(s.cards)) * CutCardPenetration)
}

// NumDecks returns the number of decks in the shoe
func (s *Shoe) NumDecks() int {
	return s.numDecks
}

// Penetration returns the percentage of cards that have been dealt
func (s *Shoe) Penetration() float64 {
	totalCards := s.numDecks * 52
	cardsDealt := totalCards - len(s.cards)
	return float64(cardsDealt) / float64(totalCards) * 100
}

// String returns a string representation of the shoe
func (s *Shoe) String() string {
	return fmt.Sprintf("Shoe: %d decks, %d cards remaining (%.1f%% penetration)",
		s.numDecks, len(s.cards), s.Penetration())
}
