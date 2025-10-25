package blackjack

import (
	"fmt"
	"strings"
	"time"

	"github.com/rbrabson/cards"
)

// ActionType represents the type of action taken on a hand
type ActionType string

const (
	ActionDeal      ActionType = "deal"
	ActionHit       ActionType = "hit"
	ActionStand     ActionType = "stand"
	ActionDouble    ActionType = "double"
	ActionSplit     ActionType = "split"
	ActionSurrender ActionType = "surrender"
)

// Action represents an action taken on a hand
type Action struct {
	Type      ActionType  `json:"type"`
	Card      *cards.Card `json:"card,omitempty"` // Card involved (for deal/hit)
	Timestamp time.Time   `json:"timestamp"`
	Details   string      `json:"details,omitempty"` // Additional details about the action
}

// Hand represents a hand of cards in blackjack
type Hand struct {
	cards    []cards.Card // cards are the game cards in the hand
	isSplit  bool         // Whether this hand came from a split
	isActive bool         // Whether this hand is still being played
	isStood  bool         // Whether the player has stood on this hand
	actions  []Action     // All actions taken on this hand
	bet      int          // The bet amount for this specific hand
	winnings int          // The winnings for this specific hand (can be negative for losses)
	player   *Player      // The player who owns this hand (nil for dealer)
}

// NewDealerHand creates a new dealer hand without a chip manager
func NewDealerHand() *Hand {
	return NewHand(nil)
}

// NewHand creates a new empty hand
func NewHand(player *Player) *Hand {
	return &Hand{
		cards:    make([]cards.Card, 0, 2),
		isSplit:  false,
		isActive: true,
		isStood:  false,
		actions:  make([]Action, 0, 1),
		bet:      0,
		winnings: 0,
		player:   player,
	}
}

// NewSplitHand creates a new hand from a split with the initial card
func NewSplitHand(card cards.Card, player *Player) *Hand {
	h := NewHand(player)
	h.isSplit = true
	h.AddCardWithAction(card, ActionDeal, "split card")

	return h
}

// AddCard adds a card to the hand
func (h *Hand) AddCard(card cards.Card) {
	h.cards = append(h.cards, card)
	// Record the card as a hit action (dealing will be tracked separately)
	h.recordAction(ActionHit, &card, "")
}

// AddCardWithAction adds a card to the hand and records the specific action
func (h *Hand) AddCardWithAction(card cards.Card, actionType ActionType, details string) {
	h.cards = append(h.cards, card)
	h.recordAction(actionType, &card, details)
}

// recordAction records an action taken on this hand
func (h *Hand) recordAction(actionType ActionType, card *cards.Card, details string) {
	action := Action{
		Type:      actionType,
		Card:      card,
		Timestamp: time.Now(),
		Details:   details,
	}
	h.actions = append(h.actions, action)
}

// RecordAction records an action without a card (like stand, surrender)
func (h *Hand) RecordAction(actionType ActionType, details string) {
	h.recordAction(actionType, nil, details)
}

// Actions returns a copy of all actions taken on this hand
func (h *Hand) Actions() []Action {
	result := make([]Action, len(h.actions))
	copy(result, h.actions)
	return result
}

// ActionSummary returns a string summary of all actions taken on this hand
func (h *Hand) ActionSummary() string {
	if len(h.actions) == 0 {
		return "No actions"
	}

	var summary strings.Builder
	for i, action := range h.actions {
		if i > 0 {
			summary.WriteString(", ")
		}

		switch action.Type {
		case ActionDeal:
			if action.Card != nil {
				summary.WriteString(fmt.Sprintf("dealt %s", action.Card))
			} else {
				summary.WriteString("dealt")
			}
		case ActionHit:
			if action.Card != nil {
				summary.WriteString(fmt.Sprintf("hit %s", action.Card))
			} else {
				summary.WriteString("hit")
			}
		case ActionStand:
			summary.WriteString("stand")
		case ActionDouble:
			if action.Card != nil {
				summary.WriteString(fmt.Sprintf("double %s", action.Card))
			} else {
				summary.WriteString("double")
			}
		case ActionSplit:
			summary.WriteString("split")
		case ActionSurrender:
			summary.WriteString("surrender")
		default:
			summary.WriteString(string(action.Type))
		}

		if action.Details != "" {
			summary.WriteString(" (")
			summary.WriteString(action.Details)
			summary.WriteString(")")
		}
	}

	return summary.String()
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
	return len(h.cards) == 2 && h.Value() == 21 && !h.IsSplit()
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
	h.bet = 0
	h.winnings = 0
}

// Bet returns the bet amount for this hand
func (h *Hand) Bet() int {
	return h.bet
}

// SetBet sets the bet amount for this hand
func (h *Hand) SetBet(amount int) {
	h.bet = amount
}

// Winnings returns the winnings for this hand (can be negative for losses)
func (h *Hand) Winnings() int {
	return h.winnings
}

// SetWinnings sets the winnings for this hand
func (h *Hand) SetWinnings(amount int) {
	h.winnings = amount
}

// AddWinnings adds to the winnings for this hand
func (h *Hand) AddWinnings(amount int) {
	h.winnings += amount
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
	h.RecordAction(ActionStand, "")
}

// CanDoubleDown returns true if the hand can be doubled down
func (h *Hand) CanDoubleDown() bool {
	return len(h.cards) == 2 && h.player.chipManager != nil && h.player.chipManager.HasEnoughChips(h.bet)
}

// DoubleDown performs the double down action on the hand
func (h *Hand) DoubleDown(card cards.Card) error {
	if !h.CanDoubleDown() {
		return fmt.Errorf("cannot double down on this hand")
	}

	// Deduct additional bet from chip manager
	err := h.player.chipManager.DeductChips(h.bet)
	if err != nil {
		return fmt.Errorf("failed to deduct chips for double down: %v", err)
	}

	// Double the bet
	h.bet *= 2

	// Add the card to the hand
	h.AddCardWithAction(card, ActionDouble, "double down")

	// Mark hand as stood after double down
	h.Stand()

	h.RecordAction(ActionDouble, fmt.Sprintf("bet increased from %d to %d", h.bet/2, h.bet))

	return nil
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
	newHand := NewSplitHand(secondCard, h.player)

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
