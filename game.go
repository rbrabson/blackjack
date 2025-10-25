package blackjack

import (
	"fmt"
	"log/slog"
	"strings"
)

// GameResult represents the outcome of a hand
type GameResult int

const (
	_               GameResult = iota
	PlayerWin                  // PlayerWin reprepsents a win for the player
	DealerWin                  // DealerWin represents a win for the dealer
	Push                       // Push represents a tie
	PlayerBlackjack            // PlayerBlackjack represents a player blackjack
	DealerBlackjack            // DealerBlackjack represents a dealer blackjack
)

// String returns a string representation of the game result
func (gr GameResult) String() string {
	switch gr {
	case PlayerWin:
		return "Player Wins"
	case DealerWin:
		return "Dealer Wins"
	case Push:
		return "Push (Tie)"
	case PlayerBlackjack:
		return "Player Blackjack!"
	case DealerBlackjack:
		return "Dealer Blackjack!"
	default:
		return "Unknown"
	}
}

// Game represents the main game
type Game struct {
	dealer  *Dealer   // dealer is the game dealer
	players []*Player // players are the game players
	shoe    *Shoe     // shoe are the cards used in the game
	round   int       // round is the current round number
}

// New creates a new blackjack game
func New(numDecks int) *Game {
	return &Game{
		dealer:  NewDealer(),
		players: make([]*Player, 0, 1),
		shoe:    NewShoe(numDecks),
		round:   0,
	}
}

// AddPlayer adds a player to the game
func (bg *Game) AddPlayer(name string, options ...Option) {
	player := NewPlayer(name, options...)
	bg.players = append(bg.players, player)
}

// GetPlayer returns a player by name
func (bg *Game) GetPlayer(name string) *Player {
	for _, player := range bg.players {
		if player.Name() == name {
			return player
		}
	}
	return nil
}

// RemovePlayer removes a player from the game
func (bg *Game) RemovePlayer(name string) bool {
	for i, player := range bg.players {
		if player.Name() == name {
			bg.players = append(bg.players[:i], bg.players[i+1:]...)
			return true
		}
	}
	return false
}

// Players returns a copy of the players slice
func (bg *Game) Players() []*Player {
	result := make([]*Player, len(bg.players))
	copy(result, bg.players)
	return result
}

// Dealer returns the dealer
func (bg *Game) Dealer() *Dealer {
	return bg.dealer
}

// Shoe returns the shoe
func (bg *Game) Shoe() *Shoe {
	return bg.shoe
}

// Round returns the current round number
func (bg *Game) Round() int {
	return bg.round
}

// DealCard deals a card from the shoe
func (bg *Game) DealCard() error {
	if bg.shoe.IsEmpty() {
		return fmt.Errorf("shoe is empty")
	}

	if bg.shoe.NeedsReshuffle() {
		slog.Debug("Reshuffling blackjack shoe...")
		bg.shoe.Reshuffle()
	}

	return nil
}

// StartNewRound starts a new round of blackjack
func (bg *Game) StartNewRound() error {
	bg.round++

	// Clear all hands
	bg.dealer.ClearHand()
	for _, player := range bg.players {
		player.ClearHand()
		player.SetActive(true)
	}

	// Check if we need to reshuffle
	if bg.shoe.NeedsReshuffle() {
		slog.Debug("Reshuffling blackjack shoe...")
		bg.shoe.Reshuffle()
	}

	return nil
}

// DealInitialCards deals two cards to each player and dealer
func (bg *Game) DealInitialCards() error {
	// Deal first card to each player
	for _, player := range bg.players {
		if player.IsActive() {
			card, err := bg.shoe.Draw()
			if err != nil {
				return fmt.Errorf("failed to deal card to %s: %w", player.Name(), err)
			}
			player.CurrentHand().DealCard(card)
		}
	}

	// Deal first card to dealer
	card, err := bg.shoe.Draw()
	if err != nil {
		return fmt.Errorf("failed to deal card to dealer: %w", err)
	}
	bg.dealer.DealCard(card)

	// Deal second card to each player
	for _, player := range bg.players {
		if player.IsActive() {
			card, err := bg.shoe.Draw()
			if err != nil {
				return fmt.Errorf("failed to deal card to %s: %w", player.Name(), err)
			}
			player.CurrentHand().DealCard(card)
		}
	}

	// Deal second card to dealer (hole card)
	card, err = bg.shoe.Draw()
	if err != nil {
		return fmt.Errorf("failed to deal hole card to dealer: %w", err)
	}
	bg.dealer.DealCard(card)

	return nil
}

// PlayerHit deals a card to a specific player
func (bg *Game) PlayerHit(playerName string) error {
	player := bg.GetPlayer(playerName)
	if player == nil {
		return fmt.Errorf("player %s not found", playerName)
	}

	if !player.IsActive() {
		return fmt.Errorf("player %s is not active", playerName)
	}

	if player.IsStanding() {
		return fmt.Errorf("player %s is already standing", playerName)
	}

	card, err := bg.shoe.Draw()
	if err != nil {
		return fmt.Errorf("failed to deal card: %w", err)
	}

	player.CurrentHand().Hit(card)
	return nil
}

// PlayerDoubleDownHit deals a card to a specific player as part of a double down
func (bg *Game) PlayerDoubleDownHit(playerName string) error {
	player := bg.GetPlayer(playerName)
	if player == nil {
		return fmt.Errorf("player %s not found", playerName)
	}

	if !player.IsActive() {
		return fmt.Errorf("player %s is not active", playerName)
	}

	card, err := bg.shoe.Draw()
	if err != nil {
		return fmt.Errorf("failed to deal card: %w", err)
	}

	player.CurrentHand().DoubleDownHit(card)
	return nil
}

// PlayerSplit handles a player splitting their hand
func (bg *Game) PlayerSplit(playerName string) error {
	player := bg.GetPlayer(playerName)
	if player == nil {
		return fmt.Errorf("player %s not found", playerName)
	}

	if !player.IsActive() {
		return fmt.Errorf("player %s is not active", playerName)
	}

	if !player.CanSplit(player.CurrentHand()) {
		return fmt.Errorf("player %s cannot split", playerName)
	}

	// Split the hand
	err := player.Split(player.CurrentHand())
	if err != nil {
		return fmt.Errorf("failed to split hand: %w", err)
	}

	// Deal a second card to each of the split hands
	hands := player.Hands()
	for i := len(hands) - 2; i < len(hands); i++ { // Last two hands are the split hands
		card, err := bg.shoe.Draw()
		if err != nil {
			return fmt.Errorf("failed to deal card to split hand: %w", err)
		}

		// Temporarily set the hand to add the card
		originalHandIdx := player.GetCurrentHandIndex()
		player.SetCurrentHandIndex(i)
		player.CurrentHand().Hit(card)
		player.SetCurrentHandIndex(originalHandIdx)
	}

	return nil
}

// PlayerStand handles a player standing on their current hand
func (bg *Game) PlayerStand(playerName string) error {
	player := bg.GetPlayer(playerName)
	if player == nil {
		return fmt.Errorf("player %s not found", playerName)
	}

	if !player.IsActive() {
		return fmt.Errorf("player %s is not active", playerName)
	}

	// Stand on current hand
	player.CurrentHand().Stand()

	// Move to next active hand if available
	if !player.MoveToNextActiveHand() {
		// No more active hands, player is done
		player.SetActive(false)
	}

	return nil
}

// PlayerSurrender handles a player surrendering their current hand
func (bg *Game) PlayerSurrender(playerName string) error {
	player := bg.GetPlayer(playerName)
	hand := player.CurrentHand()
	if player == nil {
		return fmt.Errorf("player %s not found", playerName)
	}

	if !player.IsActive() {
		return fmt.Errorf("player %s is not active", playerName)
	}

	if !hand.CanSurrender() {
		return fmt.Errorf("player %s cannot surrender at this time", playerName)
	}

	// Surrender the current hand
	hand.Surrender()

	// Move to next active hand if available
	if !player.MoveToNextActiveHand() {
		// No more active hands, player is done
		player.SetActive(false)
	}

	return nil
}

// DealerPlay handles the dealer's turn according to blackjack rules
func (bg *Game) DealerPlay() error {
	for bg.dealer.ShouldHit() {
		card, err := bg.shoe.Draw()
		if err != nil {
			return fmt.Errorf("failed to deal card to dealer: %w", err)
		}
		bg.dealer.Hit(card)
	}
	// Record that dealer is standing
	bg.dealer.Stand()
	return nil
}

// EvaluateHand determines the result of a player's hand against the dealer
func (bg *Game) EvaluateHand(playerHand *Hand) GameResult {
	dealerHand := bg.dealer.Hand()

	playerBlackjack := playerHand.IsBlackjack()
	dealerBlackjack := dealerHand.IsBlackjack()
	playerValue := playerHand.Value()
	dealerValue := dealerHand.Value()

	switch {
	case playerBlackjack && dealerBlackjack:
		return Push
	case playerBlackjack:
		return PlayerBlackjack
	case dealerBlackjack:
		return DealerBlackjack
	case playerHand.IsBusted():
		return DealerWin
	case dealerHand.IsBusted():
		return PlayerWin
	case playerValue > dealerValue:
		return PlayerWin
	case dealerValue > playerValue:
		return DealerWin
	default:
		return Push
	}
}

// PayoutResults handles payouts for all players
func (bg *Game) PayoutResults() {
	for _, player := range bg.players {
		for _, hand := range player.Hands() {
			// Skip hands with no bet
			if hand.Bet() == 0 {
				continue
			}

			result := bg.EvaluateHand(hand)

			switch result {
			case PlayerWin:
				hand.WinBet(1.0) // 1:1 payout
			case PlayerBlackjack:
				hand.WinBet(1.5) // 1.5:1 payout for blackjack
			case Push:
				hand.PushBet() // Return bet
			case DealerWin, DealerBlackjack:
				hand.LoseBet() // Lose bet
			}
		}
	}
}

// GetGameStatus returns a string representation of the current game state
func (bg *Game) GetGameStatus(showDealerHole bool) string {
	var status strings.Builder

	status.WriteString(fmt.Sprintf("=== Round %d ===\n", bg.round))
	status.WriteString(fmt.Sprintf("%s\n", bg.shoe.String()))
	status.WriteString("\n")

	// Show dealer
	if showDealerHole {
		status.WriteString(fmt.Sprintf("%s\n", bg.dealer.String()))
	} else {
		status.WriteString(fmt.Sprintf("%s\n", bg.dealer.StringHidden()))
	}
	status.WriteString("\n")

	// Show players
	for _, player := range bg.players {
		status.WriteString(fmt.Sprintf("%s\n", player.String()))
	}

	return status.String()
}

// IsRoundComplete returns true if all players have finished their hands
func (bg *Game) IsRoundComplete() bool {
	for _, player := range bg.players {
		if player.IsActive() && !player.IsStanding() {
			return false
		}
	}
	return true
}

// GetActivePlayer returns the first active player who hasn't finished their hand
func (bg *Game) GetActivePlayer() *Player {
	for _, player := range bg.players {
		if player.IsActive() && !player.IsStanding() {
			return player
		}
	}
	return nil
}
