package blackjack

import (
	"fmt"

	"github.com/rbrabson/cards"
)

// Player represents a blackjack player
type Player struct {
	name   string
	hand   *Hand
	chips  int
	bet    int
	active bool
}

// NewPlayer creates a new player with the given name and initial chips
func NewPlayer(name string, chips int) *Player {
	return &Player{
		name:   name,
		hand:   NewHand(),
		chips:  chips,
		bet:    0,
		active: true,
	}
}

// Name returns the player's name
func (p *Player) Name() string {
	return p.name
}

// Hand returns the player's hand
func (p *Player) Hand() *Hand {
	return p.hand
}

// Chips returns the player's current chip count
func (p *Player) Chips() int {
	return p.chips
}

// Bet returns the player's current bet
func (p *Player) Bet() int {
	return p.bet
}

// IsActive returns whether the player is still active in the game
func (p *Player) IsActive() bool {
	return p.active
}

// SetActive sets the player's active status
func (p *Player) SetActive(active bool) {
	p.active = active
}

// PlaceBet places a bet for the player
func (p *Player) PlaceBet(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("bet must be positive")
	}
	if amount > p.chips {
		return fmt.Errorf("insufficient chips: have %d, need %d", p.chips, amount)
	}

	p.bet = amount
	p.chips -= amount
	return nil
}

// WinBet adds winnings to the player's chips
func (p *Player) WinBet(multiplier float64) {
	winnings := int(float64(p.bet) * multiplier)
	p.chips += p.bet + winnings
	p.bet = 0
}

// LoseBet removes the player's bet (already deducted when placed)
func (p *Player) LoseBet() {
	p.bet = 0
}

// PushBet returns the bet to the player (tie)
func (p *Player) PushBet() {
	p.chips += p.bet
	p.bet = 0
}

// Hit adds a card to the player's hand
func (p *Player) Hit(card cards.Card) {
	p.hand.AddCard(card)
}

// CanDoubleDown returns true if the player can double down
func (p *Player) CanDoubleDown() bool {
	return p.hand.Count() == 2 && p.chips >= p.bet
}

// DoubleDown doubles the player's bet and they get exactly one more card
func (p *Player) DoubleDown() error {
	if !p.CanDoubleDown() {
		return fmt.Errorf("cannot double down")
	}

	p.chips -= p.bet
	p.bet *= 2
	return nil
}

// CanSplit returns true if the player can split their hand
func (p *Player) CanSplit() bool {
	if p.hand.Count() != 2 {
		return false
	}

	cards := p.hand.Cards()
	return cards[0].Rank == cards[1].Rank && p.chips >= p.bet
}

// ClearHand clears the player's hand for a new round
func (p *Player) ClearHand() {
	p.hand.Clear()
}

// String returns a string representation of the player
func (p *Player) String() string {
	status := "active"
	if !p.active {
		status = "inactive"
	}

	return fmt.Sprintf("%s (Chips: %d, Bet: %d, %s): %s",
		p.name, p.chips, p.bet, status, p.hand.String())
}

// IsStanding returns true if the player should stand (busted, blackjack, or inactive)
func (p *Player) IsStanding() bool {
	return !p.active || p.hand.IsBusted() || p.hand.IsBlackjack()
}
