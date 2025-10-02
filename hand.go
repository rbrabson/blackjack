package blackjack

import (
	"fmt"
	"strings"

	"github.com/rbrabson/cards"
)

// Hand represents a hand of cards in blackjack
type Hand struct {
	cards    []cards.Card // cards are the game cards in the hand
	isSplit  bool         // Whether this hand came from a split
	isActive bool         // Whether this hand is still being played
	isStood  bool         // Whether the player has stood on this hand
}

// NewHand creates a new empty hand
func NewHand() *Hand {
	return &Hand{
		cards:    make([]cards.Card, 0, 2),
		isSplit:  false,
		isActive: true,
		isStood:  false,
	}
}

// NewSplitHand creates a new hand from a split with the initial card
func NewSplitHand(card cards.Card) *Hand {
	return &Hand{
		cards:    []cards.Card{card},
		isSplit:  true,
		isActive: true,
		isStood:  false,
	}
}

// AddCard adds a card to the hand
func (h *Hand) AddCard(card cards.Card) {
	h.cards = append(h.cards, card)
}

// Cards returns a copy of the cards in the hand
func (h *Hand) Cards() []cards.Card {
	result := make([]cards.Card, len(h.cards))
	copy(result, h.cards)
	return result
}

// Value calculates the blackjack value of the hand
func (h *Hand) Value() int {
	value := 0
	aces := 0

	for _, card := range h.cards {
		rank := card.Rank
		switch rank {
		case cards.Jack, cards.Queen, cards.King:
			value += 10
		case cards.Ace:
			aces++
			value += 11
		default:
			value += int(rank)
		}
	}

	// Adjust for aces if value is over 21
	for aces > 0 && value > 21 {
		value -= 10
		aces--
	}

	return value
}

// IsBusted returns true if the hand value is over 21
func (h *Hand) IsBusted() bool {
	return h.Value() > 21
}

// IsBlackjack returns true if the hand is a natural blackjack (21 with 2 cards)
func (h *Hand) IsBlackjack() bool {
	return len(h.cards) == 2 && h.Value() == 21 && !h.isSplit
}

// IsSoft returns true if the hand contains an ace counted as 11
func (h *Hand) IsSoft() bool {
	value := 0
	hasAce := false

	for _, card := range h.cards {
		rank := card.Rank
		switch rank {
		case cards.Jack, cards.Queen, cards.King:
			value += 10
		case cards.Ace:
			hasAce = true
			value += 11
		default:
			value += int(rank)
		}
	}

	return hasAce && value <= 21
}

// IsSplit returns true if this hand was created by a split.
func (h *Hand) IsSplit() bool {
	return h.isSplit
}

// Count returns the number of cards in the hand
func (h *Hand) Count() int {
	return len(h.cards)
}

// Clear removes all cards from the hand
func (h *Hand) Clear() {
	h.cards = h.cards[:0]
	h.isSplit = false
	h.isActive = true
	h.isStood = false
}

// IsActive returns true if this hand is still being played
func (h *Hand) IsActive() bool {
	return h.isActive
}

// SetActive sets the active status of the hand
func (h *Hand) SetActive(active bool) {
	h.isActive = active
}

// IsStood returns true if the player has stood on this hand
func (h *Hand) IsStood() bool {
	return h.isStood
}

// Stand marks the hand as stood
func (h *Hand) Stand() {
	h.isStood = true
	h.isActive = false
}

// CanSplit returns true if the hand can be split (two cards of same rank)
func (h *Hand) CanSplit() bool {
	if len(h.cards) != 2 {
		return false
	}
	return h.cards[0].Rank == h.cards[1].Rank
}

// SplitHand splits the hand into two hands, returning the second hand
func (h *Hand) SplitHand() *Hand {
	if !h.CanSplit() {
		return nil
	}

	// Take the second card for the new hand
	secondCard := h.cards[1]
	h.cards = h.cards[:1]

	// Mark this hand as split
	h.isSplit = true

	// Create new hand with the second card
	newHand := NewSplitHand(secondCard)

	return newHand
}

// String returns a string representation of the hand
func (h *Hand) String() string {
	if len(h.cards) == 0 {
		return "Empty hand"
	}

	var cardStrings []string
	for _, card := range h.cards {
		cardStrings = append(cardStrings, card.String())
	}

	splitText := ""
	if h.isSplit {
		splitText = " (Split)"
	}

	return fmt.Sprintf("[%s] (Value: %d)%s", strings.Join(cardStrings, ", "), h.Value(), splitText)
}

// StringHidden returns a string representation with the first card hidden (for dealer)
func (h *Hand) StringHidden() string {
	if len(h.cards) == 0 {
		return "Empty hand"
	}
	if len(h.cards) == 1 {
		return "[Hidden]"
	}

	var cardStrings []string
	cardStrings = append(cardStrings, "Hidden")
	for i := 1; i < len(h.cards); i++ {
		cardStrings = append(cardStrings, h.cards[i].String())
	}

	// Calculate visible value (excluding first card)
	visibleValue := 0
	aces := 0
	for i := 1; i < len(h.cards); i++ {
		rank := h.cards[i].Rank
		switch rank {
		case cards.Jack, cards.Queen, cards.King:
			visibleValue += 10
		case cards.Ace:
			aces++
			visibleValue += 11
		default:
			visibleValue += int(rank)
		}
	}

	// Adjust for aces
	for aces > 0 && visibleValue > 21 {
		visibleValue -= 10
		aces--
	}

	return fmt.Sprintf("[%s] (Visible Value: %d)", strings.Join(cardStrings, ", "), visibleValue)
}
