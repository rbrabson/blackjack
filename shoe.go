package blackjack

import (
	"fmt"

	"github.com/rbrabson/cards"
)

// BlackjackShoe wraps the cards.Shoe with blackjack-specific functionality
type BlackjackShoe struct {
	shoe     cards.Shoe // shoe is the set of cards to be dealt
	numDecks int        // numDecdks is the number of decks in the shoe
	cutCard  int        // Position where cut card is placed (reshuffle point)
}

// NewBlackjackShoe creates a new blackjack shoe with the specified number of decks
func NewBlackjackShoe(numDecks int) *BlackjackShoe {
	if numDecks < 1 {
		numDecks = 1
	}

	shoe := cards.NewShoe(numDecks)
	shoe.Shuffle()

	// Place cut card at roughly 75% through the shoe (common casino practice)
	cutCard := int(float64(len(shoe)) * 0.75)

	return &BlackjackShoe{
		shoe:     shoe,
		numDecks: numDecks,
		cutCard:  cutCard,
	}
}

// Draw deals a card from the shoe
func (bs *BlackjackShoe) Draw() (cards.Card, error) {
	if bs.IsEmpty() {
		return cards.Card{}, fmt.Errorf("shoe is empty")
	}

	return bs.shoe.Draw(), nil
}

// IsEmpty returns true if the shoe is empty
func (bs *BlackjackShoe) IsEmpty() bool {
	return len(bs.shoe) == 0
}

// NeedsReshuffle returns true if the cut card has been reached
func (bs *BlackjackShoe) NeedsReshuffle() bool {
	return len(bs.shoe) <= (bs.numDecks*52 - bs.cutCard)
}

// CardsRemaining returns the number of cards left in the shoe
func (bs *BlackjackShoe) CardsRemaining() int {
	return len(bs.shoe)
}

// Reshuffle creates a new shuffled shoe with the same number of decks
func (bs *BlackjackShoe) Reshuffle() {
	bs.shoe = cards.NewShoe(bs.numDecks)
	bs.shoe.Shuffle()

	// Reset cut card position
	bs.cutCard = int(float64(len(bs.shoe)) * 0.75)
}

// NumDecks returns the number of decks in the shoe
func (bs *BlackjackShoe) NumDecks() int {
	return bs.numDecks
}

// Penetration returns the percentage of cards that have been dealt
func (bs *BlackjackShoe) Penetration() float64 {
	totalCards := bs.numDecks * 52
	cardsDealt := totalCards - len(bs.shoe)
	return float64(cardsDealt) / float64(totalCards) * 100
}

// String returns a string representation of the shoe
func (bs *BlackjackShoe) String() string {
	return fmt.Sprintf("Shoe: %d decks, %d cards remaining (%.1f%% penetration)",
		bs.numDecks, len(bs.shoe), bs.Penetration())
}
