package blackjack

import (
	"fmt"

	"github.com/rbrabson/cards"
)

// Dealer represents the blackjack dealer
type Dealer struct {
	hand *Hand // hand is the dealer's hand
}

// NewDealer creates a new dealer
func NewDealer() *Dealer {
	return &Dealer{
		hand: NewHand(),
	}
}

// Hand returns the dealer's hand
func (d *Dealer) Hand() *Hand {
	return d.hand
}

// Hit adds a card to the dealer's hand
func (d *Dealer) Hit(card cards.Card) {
	d.hand.AddCardWithAction(card, ActionHit, "dealer hit")
}

// DealCard adds a card to the dealer's hand as part of the initial deal
func (d *Dealer) DealCard(card cards.Card) {
	d.hand.AddCardWithAction(card, ActionDeal, "initial deal")
}

// Stand marks the dealer as standing
func (d *Dealer) Stand() {
	d.hand.RecordAction(ActionStand, "dealer stands")
	d.hand.isStood = true
	d.hand.isActive = false
}

// ShouldHit returns true if the dealer should hit according to standard blackjack rules
// Dealer hits on 16 or less, stands on 17 or more (including soft 17)
func (d *Dealer) ShouldHit() bool {
	value := d.hand.Value()

	switch {
	// Always stand if busted
	case d.hand.IsBusted():
		return false
	// Stand on hard 17 or higher
	case value >= 17 && !d.hand.IsSoft():
		return false
	// Hit on soft 17 (house rule - can be changed)
	case value == 17 && d.hand.IsSoft():
		return true
	// Stand on soft 18 or higher
	case value >= 18:
		return false
	// Hit on 16 or less
	default:
		return value <= 16

	}
}

// ShowFirstCard returns the dealer's first card (face up)
func (d *Dealer) ShowFirstCard() cards.Card {
	if d.hand.Count() == 0 {
		panic("dealer has no cards")
	}
	return d.hand.Cards()[0]
}

// HasBlackjack returns true if dealer has blackjack
func (d *Dealer) HasBlackjack() bool {
	return d.hand.IsBlackjack()
}

// IsBusted returns true if dealer is busted
func (d *Dealer) IsBusted() bool {
	return d.hand.IsBusted()
}

// Value returns the dealer's hand value
func (d *Dealer) Value() int {
	return d.hand.Value()
}

// ClearHand clears the dealer's hand for a new round
func (d *Dealer) ClearHand() {
	d.hand.Clear()
}

// String returns a string representation of the dealer with both cards showing
func (d *Dealer) String() string {
	return fmt.Sprintf("Dealer: %s", d.hand.String())
}

// StringHidden returns a string representation of the dealer with hole card hidden
func (d *Dealer) StringHidden() string {
	return fmt.Sprintf("Dealer: %s", d.hand.StringHidden())
}

// RevealHoleCard shows the dealer's full hand
func (d *Dealer) RevealHoleCard() string {
	return d.String()
}
