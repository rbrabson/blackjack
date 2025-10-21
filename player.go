package blackjack

import (
	"fmt"
	"strings"

	"github.com/rbrabson/cards"
)

// Player represents a blackjack player
type Player struct {
	name           string
	hands          []Hand
	chipManager    ChipManager
	active         bool
	currentHandIdx int
}

// NewPlayer creates a new player with the given name, initial chips, and optional settings
func NewPlayer(name string, chips int, options ...Option) *Player {
	player := &Player{
		name:           name,
		hands:          []Hand{*NewHand()},
		chipManager:    NewDefaultChipManager(0),
		active:         true,
		currentHandIdx: 0,
	}
	for _, option := range options {
		option(player)
	}
	player.chipManager.SetChips(chips)
	return player
}

// Option is a function that modifies a message.
type Option func(*Player)

// Name returns the player's name
func (p *Player) Name() string {
	return p.name
}

// WithChipManager sets a custom chip manager for the player.
func WithChipManager(cm ChipManager) Option {
	return func(p *Player) {
		p.chipManager = cm
	}
}

// WithAllowedMentions sets the allowed mentions for the message.
func WithChips(chips int) Option {
	return func(p *Player) {
		p.chipManager.SetChips(chips)
	}
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

// Bet returns the player's current hand bet
func (p *Player) Bet() int {
	return p.CurrentHand().Bet()
}

// IsActive returns whether the player is still active in the game
func (p *Player) IsActive() bool {
	return p.active
}

// SetActive sets the player's active status
func (p *Player) SetActive(active bool) {
	p.active = active
}

// PlaceBet places a bet for the player's current hand
func (p *Player) PlaceBet(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("bet must be positive")
	}
	if !p.chipManager.HasEnoughChips(amount) {
		return fmt.Errorf("insufficient chips: have %d, need %d", p.chipManager.GetChips(), amount)
	}

	// Set bet on current hand and deduct from chips
	p.CurrentHand().SetBet(amount)
	return p.chipManager.DeductChips(amount)
}

// PlaceBetOnHand places a bet for a specific hand (useful for splits)
func (p *Player) PlaceBetOnHand(handIndex int, amount int) error {
	if handIndex < 0 || handIndex >= len(p.hands) {
		return fmt.Errorf("invalid hand index: %d", handIndex)
	}
	if amount <= 0 {
		return fmt.Errorf("bet must be positive")
	}
	if !p.chipManager.HasEnoughChips(amount) {
		return fmt.Errorf("insufficient chips: have %d, need %d", p.chipManager.GetChips(), amount)
	}

	p.hands[handIndex].SetBet(amount)
	return p.chipManager.DeductChips(amount)
}

// WinBet adds winnings to the player's chips for the current hand
func (p *Player) WinBet(multiplier float64) {
	hand := p.CurrentHand()
	winnings := int(float64(hand.Bet()) * multiplier)
	totalPayout := hand.Bet() + winnings
	p.chipManager.AddChips(totalPayout)
	hand.SetWinnings(winnings)
}

// WinBetOnHand adds winnings to the player's chips for a specific hand
func (p *Player) WinBetOnHand(handIndex int, multiplier float64) {
	if handIndex < 0 || handIndex >= len(p.hands) {
		return
	}
	hand := &p.hands[handIndex]
	winnings := int(float64(hand.Bet()) * multiplier)
	totalPayout := hand.Bet() + winnings
	p.chipManager.AddChips(totalPayout)
	hand.SetWinnings(winnings)
}

// LoseBet removes the player's bet for the current hand (already deducted when placed)
func (p *Player) LoseBet() {
	hand := p.CurrentHand()
	hand.SetWinnings(-hand.Bet()) // Record the loss
}

// LoseBetOnHand removes the player's bet for a specific hand
func (p *Player) LoseBetOnHand(handIndex int) {
	if handIndex < 0 || handIndex >= len(p.hands) {
		return
	}
	hand := &p.hands[handIndex]
	hand.SetWinnings(-hand.Bet()) // Record the loss
}

// PushBet returns the bet to the player for the current hand (tie)
func (p *Player) PushBet() {
	hand := p.CurrentHand()
	p.chipManager.AddChips(hand.Bet())
	hand.SetWinnings(0) // No win or loss
}

// PushBetOnHand returns the bet to the player for a specific hand
func (p *Player) PushBetOnHand(handIndex int) {
	if handIndex < 0 || handIndex >= len(p.hands) {
		return
	}
	hand := &p.hands[handIndex]
	p.chipManager.AddChips(hand.Bet())
	hand.SetWinnings(0) // No win or loss
}

// Surrender allows the player to forfeit their hand and lose half their bet
func (p *Player) Surrender() {
	currentBet := p.CurrentHand().Bet()
	halfBet := currentBet / 2
	p.chipManager.AddChips(halfBet)
	p.CurrentHand().SetWinnings(-halfBet) // Record the loss of half bet
	p.CurrentHand().SetBet(0)
	p.CurrentHand().RecordAction(ActionSurrender, fmt.Sprintf("received %d chips back", halfBet))
	p.CurrentHand().Stand()
}

// CanSurrender returns true if the player can surrender (typically only on first two cards)
func (p *Player) CanSurrender() bool {
	currentHand := p.CurrentHand()
	return len(p.Hands()) == 1 && currentHand.Count() == 2 && !currentHand.IsStood() && !currentHand.IsBusted()
}

// Hit adds a card to the player's hand
func (p *Player) Hit(card cards.Card) {
	// Use AddCardWithAction to specify this is a hit
	p.CurrentHand().AddCardWithAction(card, ActionHit, "player hit")
}

// DealCard adds a card to the player's hand as part of the initial deal
func (p *Player) DealCard(card cards.Card) {
	p.CurrentHand().AddCardWithAction(card, ActionDeal, "initial deal")
}

// DoubleDownHit adds a card to the player's hand as part of a double down
func (p *Player) DoubleDownHit(card cards.Card) {
	p.CurrentHand().AddCardWithAction(card, ActionDouble, "double down card")
}

// CanDoubleDown returns true if the player can double down
func (p *Player) CanDoubleDown() bool {
	return p.CurrentHand().Count() == 2 && p.chipManager.HasEnoughChips(p.CurrentHand().Bet())
}

// DoubleDown doubles the player's bet and they get exactly one more card
func (p *Player) DoubleDown() error {
	if !p.CanDoubleDown() {
		return fmt.Errorf("cannot double down")
	}

	currentBet := p.CurrentHand().Bet()
	err := p.chipManager.DeductChips(currentBet)
	if err != nil {
		return err
	}
	newBet := currentBet * 2
	p.CurrentHand().SetBet(newBet)
	p.CurrentHand().RecordAction(ActionDouble, fmt.Sprintf("bet increased from %d to %d", currentBet, newBet))
	return nil
}

// Split splits the player's hand into two hands
func (p *Player) Split() error {
	if !p.CanSplit() {
		return fmt.Errorf("cannot split")
	}

	currentHand := p.CurrentHand()

	// Record split action before splitting
	currentHand.RecordAction(ActionSplit, fmt.Sprintf("split into %d hands", len(p.hands)+1))

	// Use the Hand's SplitHand method to get the new hand
	newHand := currentHand.SplitHand()
	if newHand == nil {
		return fmt.Errorf("split failed")
	}

	// Set the same bet on the new hand before adding to slice
	currentBet := currentHand.Bet()
	newHand.SetBet(currentBet)

	// Record split action on the new hand too
	newHand.RecordAction(ActionSplit, "created from split")

	// Add the new hand to the player's hands
	p.hands = append(p.hands, *newHand)

	// Deduct from chips for the new hand's bet
	err := p.chipManager.DeductChips(currentBet)
	return err
}

// CanSplit returns true if the player can split their hand
func (p *Player) CanSplit() bool {
	// Can only split if we have enough chips, the hand can be split, and we have fewer than 4 hands (maximum allowed)
	return len(p.hands) < 4 && p.CurrentHand().CanSplit() && p.chipManager.HasEnoughChips(p.CurrentHand().Bet())
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
			p.name, p.chipManager.GetChips(), p.hands[0].Bet(), status, p.hands[0].String())
	} else {
		// Multiple hands (splits) - show total bet across all hands
		totalBet := 0
		for _, hand := range p.hands {
			totalBet += hand.Bet()
		}
		handStrings := make([]string, len(p.hands))
		for i, hand := range p.hands {
			current := ""
			if i == p.currentHandIdx {
				current = " *CURRENT*"
			}
			handStrings[i] = fmt.Sprintf("Hand %d (Bet: %d): %s%s", i+1, hand.Bet(), hand.String(), current)
		}
		return fmt.Sprintf("%s (Chips: %d, Total Bet: %d, %s):\n  %s",
			p.name, p.chipManager.GetChips(), totalBet, status, strings.Join(handStrings, "\n  "))
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
