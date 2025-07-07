package blackjack

import "github.com/rbrabson/blackjack/deck"

type Hand struct {
	Cards []deck.Card
}

func NewHand() *Hand {
	return &Hand{
		Cards: []deck.Card{},
	}
}

func (h *Hand) AddCard(card deck.Card) {
	h.Cards = append(h.Cards, card)
}

func (h *Hand) Clear() {
	h.Cards = []deck.Card{}
}

func (h *Hand) Split() (*Hand, *Hand) {
	if len(h.Cards) < 2 {
		return NewHand(), NewHand()
	}

	hand1 := &Hand{Cards: []deck.Card{h.Cards[0]}}
	hand2 := &Hand{Cards: []deck.Card{h.Cards[1]}}

	// Clear the original hand
	h.Clear()

	return hand1, hand2
}
