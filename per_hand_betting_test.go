package blackjack

import (
	"testing"

	"github.com/rbrabson/cards"
)

// TestPerHandBetting tests basic per-hand betting functionality
func TestPerHandBetting(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Test initial state
	if player.Bet() != 0 {
		t.Error("Initial bet should be 0")
	}
	if player.CurrentHand().Bet() != 0 {
		t.Error("Initial hand bet should be 0")
	}
	if player.CurrentHand().Winnings() != 0 {
		t.Error("Initial hand winnings should be 0")
	}

	// Place a bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Check bet was set correctly
	if player.Bet() != 100 {
		t.Errorf("Expected bet of 100, got %d", player.Bet())
	}
	if player.CurrentHand().Bet() != 100 {
		t.Errorf("Expected hand bet of 100, got %d", player.CurrentHand().Bet())
	}
	if player.Chips() != 900 {
		t.Errorf("Expected 900 chips after betting, got %d", player.Chips())
	}

	// Test win
	player.WinBet(1.0) // 1:1 payout
	if player.CurrentHand().Bet() == 0 {
		t.Error("Bet shouldn't be cleared after win")
	}
	if player.CurrentHand().Winnings() != 100 {
		t.Errorf("Expected winnings of 100, got %d", player.CurrentHand().Winnings())
	}
	if player.Chips() != 1100 {
		t.Errorf("Expected 1100 chips after win, got %d", player.Chips())
	}
}

// TestPerHandBettingWithSplit tests betting with split hands
func TestPerHandBettingWithSplit(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up for split
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Eight}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Eight}

	player.DealCard(card1)
	player.DealCard(card2)

	// Place initial bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	initialChips := player.Chips() // Should be 900

	// Split
	err = player.Split()
	if err != nil {
		t.Fatalf("Failed to split: %v", err)
	}

	// Check that we now have 2 hands with correct bets
	if len(player.Hands()) != 2 {
		t.Errorf("Expected 2 hands after split, got %d", len(player.Hands()))
	}

	// Both hands should have the same bet
	if player.hands[0].Bet() != 100 {
		t.Errorf("Expected first hand bet of 100, got %d", player.hands[0].Bet())
	}
	if player.hands[1].Bet() != 100 {
		t.Errorf("Expected second hand bet of 100, got %d", player.hands[1].Bet())
	}

	// Chips should be reduced by the additional bet
	expectedChips := initialChips - 100 // Additional 100 for the split
	if player.Chips() != expectedChips {
		t.Errorf("Expected %d chips after split, got %d", expectedChips, player.Chips())
	}

	// Test individual hand betting methods
	player.WinBetOnHand(0, 1.0) // First hand wins 1:1
	player.LoseBetOnHand(1)     // Second hand loses

	// Check results
	if player.hands[0].Bet() == 0 {
		t.Error("First hand bet shouldn't be cleared after win")
	}
	if player.hands[0].Winnings() != 100 {
		t.Errorf("Expected first hand winnings of 100, got %d", player.hands[0].Winnings())
	}

	if player.hands[1].Bet() == 0 {
		t.Error("Second hand bet should be cleared after loss")
	}
	if player.hands[1].Winnings() != -100 {
		t.Errorf("Expected second hand winnings of -100, got %d", player.hands[1].Winnings())
	}

	// First hand should return bet + winnings = 200 chips
	expectedFinalChips := expectedChips + 200
	if player.Chips() != expectedFinalChips {
		t.Errorf("Expected %d chips after payouts, got %d", expectedFinalChips, player.Chips())
	}
}

// TestPerHandBettingWithDoubleDown tests betting with double down
func TestPerHandBettingWithDoubleDown(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up for double down
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.DealCard(card1)
	player.DealCard(card2)

	// Place initial bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Double down
	err = player.DoubleDown()
	if err != nil {
		t.Fatalf("Failed to double down: %v", err)
	}

	// Check bet was doubled
	if player.CurrentHand().Bet() != 200 {
		t.Errorf("Expected doubled bet of 200, got %d", player.CurrentHand().Bet())
	}

	// Check chips were deducted for the additional bet
	if player.Chips() != 800 { // 1000 - 100 (initial) - 100 (double)
		t.Errorf("Expected 800 chips after double down, got %d", player.Chips())
	}

	// Test win with doubled bet
	player.WinBet(1.0) // 1:1 payout on doubled bet

	if player.CurrentHand().Bet() == 0 {
		t.Error("Bet shouldn't be cleared after win")
	}
	if player.CurrentHand().Winnings() != 200 {
		t.Errorf("Expected winnings of 200, got %d", player.CurrentHand().Winnings())
	}
	if player.Chips() != 1200 { // 800 + 200 (bet) + 200 (winnings)
		t.Errorf("Expected 1200 chips after win, got %d", player.Chips())
	}
}

// TestPerHandBettingSurrender tests betting with surrender
func TestPerHandBettingSurrender(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up for surrender
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.DealCard(card1)
	player.DealCard(card2)

	// Place bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Surrender
	player.Surrender()

	// Check bet was cleared and half returned
	if player.CurrentHand().Bet() != 0 {
		t.Error("Bet should be cleared after surrender")
	}
	if player.CurrentHand().Winnings() != -50 {
		t.Errorf("Expected winnings of -50 (half bet lost), got %d", player.CurrentHand().Winnings())
	}
	if player.Chips() != 950 { // 1000 - 100 + 50 (half back)
		t.Errorf("Expected 950 chips after surrender, got %d", player.Chips())
	}
}

// TestPerHandBettingPush tests betting with push (tie)
func TestPerHandBettingPush(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Place bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Push (tie)
	player.PushBet()

	// Check bet was cleared and money returned
	if player.CurrentHand().Bet() == 0 {
		t.Error("Bet shouldn't be cleared after push")
	}
	if player.CurrentHand().Winnings() != 0 {
		t.Errorf("Expected winnings of 0 (push), got %d", player.CurrentHand().Winnings())
	}
	if player.Chips() != 1000 { // Back to original amount
		t.Errorf("Expected 1000 chips after push, got %d", player.Chips())
	}
}

// TestHandClearResetsFields tests that clearing a hand resets bet and winnings
func TestHandClearResetsFields(t *testing.T) {
	hand := NewHand()

	// Set some values
	hand.SetBet(100)
	hand.SetWinnings(50)

	// Verify they're set
	if hand.Bet() != 100 || hand.Winnings() != 50 {
		t.Error("Bet and winnings should be set before clear")
	}

	// Clear the hand
	hand.Clear()

	// Verify they're reset
	if hand.Bet() != 0 {
		t.Errorf("Expected bet to be 0 after clear, got %d", hand.Bet())
	}
	if hand.Winnings() != 0 {
		t.Errorf("Expected winnings to be 0 after clear, got %d", hand.Winnings())
	}
}
