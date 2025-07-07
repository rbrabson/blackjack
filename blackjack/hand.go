package blackjack

import "github.com/rbrabson/blackjack/deck"

// Hand represents a player's hand in blackjack, containing a collection of cards.
type Hand struct {
	Cards []deck.Card
}

// NewHand creates a new empty hand.
func NewHand() *Hand {
	return &Hand{
		Cards: []deck.Card{},
	}
}

// AddCard adds a card to the hand.
func (h *Hand) AddCard(card deck.Card) {
	h.Cards = append(h.Cards, card)
}

// Clear removes all cards from the hand.
func (h *Hand) Clear() {
	h.Cards = []deck.Card{}
}

// Split splits the hand into two hands, returning them. The original hand is cleared.
func (h *Hand) Split() (*Hand, *Hand) {
	if len(h.Cards) < 2 {
		return NewHand(), NewHand()
	}

	hand1 := &Hand{Cards: []deck.Card{h.Cards[0]}}
	hand2 := &Hand{Cards: []deck.Card{h.Cards[1]}}

	h.Clear()

	return hand1, hand2
}
