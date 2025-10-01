package blackjack

import (
	"fmt"
	"strings"

	"github.com/rbrabson/cards"
)

// ChipManager interface defines the operations for managing player chips
type ChipManager interface {
	GetChips() int
	SetChips(amount int)
	AddChips(amount int)
	DeductChips(amount int) error
	HasEnoughChips(amount int) bool
}

// DefaultChipManager implements ChipManager with simple integer-based chip management
type DefaultChipManager struct {
	chips int
}

// NewDefaultChipManager creates a new default chip manager with the given initial amount
func NewDefaultChipManager(initialChips int) *DefaultChipManager {
	return &DefaultChipManager{chips: initialChips}
}

// GetChips returns the current chip count
func (c *DefaultChipManager) GetChips() int {
	return c.chips
}

// SetChips sets the chip count to the specified amount
func (c *DefaultChipManager) SetChips(amount int) {
	c.chips = amount
}

// AddChips adds the specified amount to the chip count
func (c *DefaultChipManager) AddChips(amount int) {
	c.chips += amount
}

// DeductChips removes the specified amount from the chip count
func (c *DefaultChipManager) DeductChips(amount int) error {
	if amount > c.chips {
		return fmt.Errorf("insufficient chips: have %d, need %d", c.chips, amount)
	}
	c.chips -= amount
	return nil
}

// HasEnoughChips returns true if there are enough chips for the specified amount
func (c *DefaultChipManager) HasEnoughChips(amount int) bool {
	return c.chips >= amount
}

// Player represents a blackjack player
type Player struct {
	name           string
	hands          []Hand
	chipManager    ChipManager
	bet            int
	active         bool
	currentHandIdx int
}

// NewPlayer creates a new player with the given name and initial chips
func NewPlayer(name string, chips int) *Player {
	return &Player{
		name:           name,
		hands:          []Hand{*NewHand()},
		chipManager:    NewDefaultChipManager(chips),
		bet:            0,
		active:         true,
		currentHandIdx: 0,
	}
}

// NewPlayerWithChipManager creates a new player with a custom chip manager
func NewPlayerWithChipManager(name string, chipManager ChipManager) *Player {
	return &Player{
		name:           name,
		hands:          []Hand{*NewHand()},
		chipManager:    chipManager,
		bet:            0,
		active:         true,
		currentHandIdx: 0,
	}
}

// Name returns the player's name
func (p *Player) Name() string {
	return p.name
}

// Hand returns all of the player's hands
func (p *Player) Hands() []Hand {
	return p.hands
}

// CurrentHand returns the player's current hand
func (p *Player) CurrentHand() *Hand {
	return &p.hands[p.currentHandIdx]
}

// NextHand moves to the next hand if available, returning true if successful
func (p *Player) NextHand() bool {
	if p.currentHandIdx+1 < len(p.hands) {
		p.currentHandIdx++
		return true
	}
	return false
}

// Chips returns the player's current chip count
func (p *Player) Chips() int {
	return p.chipManager.GetChips()
}

// AddChips adds chips to the player's account
func (p *Player) AddChips(amount int) {
	p.chipManager.AddChips(amount)
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
	if !p.chipManager.HasEnoughChips(amount) {
		return fmt.Errorf("insufficient chips: have %d, need %d", p.chipManager.GetChips(), amount)
	}

	p.bet = amount
	return p.chipManager.DeductChips(amount)
}

// WinBet adds winnings to the player's chips
func (p *Player) WinBet(multiplier float64) {
	winnings := int(float64(p.bet) * multiplier)
	p.chipManager.AddChips(p.bet + winnings)
	p.bet = 0
}

// LoseBet removes the player's bet (already deducted when placed)
func (p *Player) LoseBet() {
	p.bet = 0
}

// PushBet returns the bet to the player (tie)
func (p *Player) PushBet() {
	p.chipManager.AddChips(p.bet)
	p.bet = 0
}

// Hit adds a card to the player's hand
func (p *Player) Hit(card cards.Card) {
	p.hands[p.currentHandIdx].AddCard(card)
}

// CanDoubleDown returns true if the player can double down
func (p *Player) CanDoubleDown() bool {
	return p.hands[p.currentHandIdx].Count() == 2 && p.chipManager.HasEnoughChips(p.bet)
}

// DoubleDown doubles the player's bet and they get exactly one more card
func (p *Player) DoubleDown() error {
	if !p.CanDoubleDown() {
		return fmt.Errorf("cannot double down")
	}

	err := p.chipManager.DeductChips(p.bet)
	if err != nil {
		return err
	}
	p.bet *= 2
	return nil
}

// Split splits the player's hand into two hands
func (p *Player) Split() error {
	if !p.CanSplit() {
		return fmt.Errorf("cannot split")
	}

	currentHand := &p.hands[p.currentHandIdx]

	// Use the Hand's SplitHand method to get the new hand
	newHand := currentHand.SplitHand()
	if newHand == nil {
		return fmt.Errorf("split failed")
	}

	// Add the new hand to the player's hands
	p.hands = append(p.hands, *newHand)

	// Deduct the bet for the new hand
	err := p.chipManager.DeductChips(p.bet)
	return err
}

// CanSplit returns true if the player can split their hand
func (p *Player) CanSplit() bool {
	// Can only split if we have enough chips and the hand can be split
	return p.hands[p.currentHandIdx].CanSplit() && p.chipManager.HasEnoughChips(p.bet)
}

// ClearHand clears all of the player's hands for a new round
func (p *Player) ClearHand() {
	// Reset to a single hand
	p.hands = []Hand{*NewHand()}
	p.currentHandIdx = 0
}

// String returns a string representation of the player
func (p *Player) String() string {
	status := "active"
	if !p.active {
		status = "inactive"
	}

	if len(p.hands) == 1 {
		// Single hand
		return fmt.Sprintf("%s (Chips: %d, Bet: %d, %s): %s",
			p.name, p.chipManager.GetChips(), p.bet, status, p.hands[0].String())
	} else {
		// Multiple hands (splits)
		handStrings := make([]string, len(p.hands))
		for i, hand := range p.hands {
			current := ""
			if i == p.currentHandIdx {
				current = " *CURRENT*"
			}
			handStrings[i] = fmt.Sprintf("Hand %d: %s%s", i+1, hand.String(), current)
		}
		return fmt.Sprintf("%s (Chips: %d, Bet: %d, %s):\n  %s",
			p.name, p.chipManager.GetChips(), p.bet, status, strings.Join(handStrings, "\n  "))
	}
}

// IsStanding returns true if the current hand should stand (busted, blackjack, or inactive)
func (p *Player) IsStanding() bool {
	if !p.active {
		return true
	}

	currentHand := &p.hands[p.currentHandIdx]
	return currentHand.IsBusted() || currentHand.IsBlackjack() || currentHand.IsStood()
}

// HasActiveHands returns true if the player has any active hands left to play
func (p *Player) HasActiveHands() bool {
	if !p.active {
		return false
	}

	for i := p.currentHandIdx; i < len(p.hands); i++ {
		hand := &p.hands[i]
		if !hand.IsBusted() && !hand.IsBlackjack() && !hand.IsStood() {
			return true
		}
	}
	return false
}

// MoveToNextActiveHand moves to the next active hand, returns true if successful
func (p *Player) MoveToNextActiveHand() bool {
	for i := p.currentHandIdx + 1; i < len(p.hands); i++ {
		if !p.hands[i].IsBusted() && !p.hands[i].IsBlackjack() && !p.hands[i].IsStood() {
			p.currentHandIdx = i
			return true
		}
	}
	return false
}

// GetAllHandValues returns the values of all hands
func (p *Player) GetAllHandValues() []int {
	values := make([]int, len(p.hands))
	for i, hand := range p.hands {
		values[i] = hand.Value()
	}
	return values
}

// GetCurrentHandIndex returns the index of the current hand
func (p *Player) GetCurrentHandIndex() int {
	return p.currentHandIdx
}

// SetCurrentHandIndex sets the current hand index (for internal use)
func (p *Player) SetCurrentHandIndex(index int) {
	if index >= 0 && index < len(p.hands) {
		p.currentHandIdx = index
	}
}
