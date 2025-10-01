package blackjack

import (
	"fmt"
	"testing"

	"github.com/rbrabson/cards"
)

func TestHandValue(t *testing.T) {
	hand := NewHand()

	// Test basic card values
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Five})
	if hand.Value() != 15 {
		t.Errorf("Expected 15, got %d", hand.Value())
	}

	// Test face cards
	hand.Clear()
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.King})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Queen})
	if hand.Value() != 20 {
		t.Errorf("Expected 20, got %d", hand.Value())
	}

	// Test blackjack
	hand.Clear()
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.King})
	if hand.Value() != 21 {
		t.Errorf("Expected 21, got %d", hand.Value())
	}
	if !hand.IsBlackjack() {
		t.Error("Should be blackjack")
	}

	// Test soft ace
	hand.Clear()
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Six})
	if hand.Value() != 17 {
		t.Errorf("Expected 17, got %d", hand.Value())
	}
	if !hand.IsSoft() {
		t.Error("Should be soft hand")
	}

	// Test ace adjustment on bust
	hand.AddCard(cards.Card{Suit: cards.Diamonds, Rank: cards.Ten})
	if hand.Value() != 17 {
		t.Errorf("Expected 17 (ace adjusted), got %d", hand.Value())
	}
	if hand.IsSoft() {
		t.Error("Should not be soft hand after adjustment")
	}
}

func TestPlayerBetting(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Test valid bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if player.Bet() != 100 {
		t.Errorf("Expected bet 100, got %d", player.Bet())
	}
	if player.Chips() != 900 {
		t.Errorf("Expected 900 chips, got %d", player.Chips())
	}

	// Test insufficient chips
	err = player.PlaceBet(1000)
	if err == nil {
		t.Error("Expected error for insufficient chips")
	}

	// Test double down
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Nine})

	if !player.CanDoubleDown() {
		t.Error("Should be able to double down")
	}

	err = player.DoubleDown()
	if err != nil {
		t.Errorf("Unexpected error during double down: %v", err)
	}
	if player.Bet() != 200 {
		t.Errorf("Expected bet 200 after double down, got %d", player.Bet())
	}
	if player.Chips() != 800 {
		t.Errorf("Expected 800 chips after double down, got %d", player.Chips())
	}
}

func TestDealerRules(t *testing.T) {
	dealer := NewDealer()

	// Test dealer hits on 16
	dealer.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	dealer.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Six})
	if !dealer.ShouldHit() {
		t.Error("Dealer should hit on 16")
	}

	// Test dealer stands on 17
	dealer.ClearHand()
	dealer.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	dealer.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Seven})
	if dealer.ShouldHit() {
		t.Error("Dealer should stand on 17")
	}

	// Test dealer hits on soft 17
	dealer.ClearHand()
	dealer.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	dealer.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Six})
	if !dealer.ShouldHit() {
		t.Error("Dealer should hit on soft 17")
	}
}

func TestGameEvaluation(t *testing.T) {
	game := New(1)
	game.AddPlayer("TestPlayer", 1000)
	player := game.GetPlayer("TestPlayer")
	player.PlaceBet(100)

	// Test player blackjack vs dealer non-blackjack
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.King})
	game.Dealer().Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	game.Dealer().Hit(cards.Card{Suit: cards.Diamonds, Rank: cards.Nine})

	result := game.EvaluateHand(player)
	if result != PlayerBlackjack {
		t.Errorf("Expected PlayerBlackjack, got %v", result)
	}

	// Test player bust
	player.ClearHand()
	game.Dealer().ClearHand()
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Ten})
	player.Hit(cards.Card{Suit: cards.Diamonds, Rank: cards.Five})
	game.Dealer().Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	game.Dealer().Hit(cards.Card{Suit: cards.Clubs, Rank: cards.Seven})

	result = game.EvaluateHand(player)
	if result != DealerWin {
		t.Errorf("Expected DealerWin (player bust), got %v", result)
	}

	// Test push
	player.ClearHand()
	game.Dealer().ClearHand()
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Nine})
	game.Dealer().Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ten})
	game.Dealer().Hit(cards.Card{Suit: cards.Clubs, Rank: cards.Nine})

	result = game.EvaluateHand(player)
	if result != Push {
		t.Errorf("Expected Push, got %v", result)
	}
}

func ExampleGame() {
	// Create a game with 1 deck for predictable testing
	game := New(1)
	game.AddPlayer("Alice", 500)

	player := game.GetPlayer("Alice")
	player.PlaceBet(50)

	fmt.Printf("Player: %s\n", player.Name())
	fmt.Printf("Chips: %d\n", player.Chips())
	fmt.Printf("Bet: %d\n", player.Bet())

	// Output:
	// Player: Alice
	// Chips: 450
	// Bet: 50
}

func TestHandSplit(t *testing.T) {
	// Test basic split functionality
	hand := NewHand()
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.Eight})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Eight})

	if !hand.CanSplit() {
		t.Error("Should be able to split pair of eights")
	}

	// Split the hand
	newHand := hand.SplitHand()
	if newHand == nil {
		t.Fatal("Split should have returned a new hand")
	}

	// Check original hand
	if hand.Count() != 1 {
		t.Errorf("Original hand should have 1 card, got %d", hand.Count())
	}
	if !hand.IsSplit() {
		t.Error("Original hand should be marked as split")
	}

	// Check new hand
	if newHand.Count() != 1 {
		t.Errorf("New hand should have 1 card, got %d", newHand.Count())
	}
	if !newHand.IsSplit() {
		t.Error("New hand should be marked as split")
	}

	// Split hands cannot have blackjack
	hand.AddCard(cards.Card{Suit: cards.Clubs, Rank: cards.Ace})
	if hand.IsBlackjack() {
		t.Error("Split hand should not be considered blackjack")
	}
}

func TestPlayerSplit(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)
	player.PlaceBet(100)

	// Add pair of kings
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.King})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.King})

	if !player.CanSplit() {
		t.Error("Should be able to split pair of kings")
	}

	// Split the hand
	err := player.Split()
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Check player state after split
	hands := player.Hands()
	if len(hands) != 2 {
		t.Errorf("Expected 2 hands after split, got %d", len(hands))
	}

	// Check chips (should be reduced by bet amount)
	if player.Chips() != 800 { // 1000 - 100 (initial bet) - 100 (split bet)
		t.Errorf("Expected 800 chips after split, got %d", player.Chips())
	}

	// Test multiple hand navigation
	if player.GetCurrentHandIndex() != 0 {
		t.Error("Should start with first hand")
	}

	if !player.NextHand() {
		t.Error("Should be able to move to next hand")
	}

	if player.GetCurrentHandIndex() != 1 {
		t.Error("Should be on second hand after NextHand")
	}
}

func TestGameSplit(t *testing.T) {
	game := New(1)
	game.AddPlayer("TestPlayer", 1000)
	player := game.GetPlayer("TestPlayer")

	// Start a new round
	game.StartNewRound()
	player.PlaceBet(100)

	// Manually set up a split scenario
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Nine})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Nine})

	// Test game split method
	err := game.PlayerSplit("TestPlayer")
	if err != nil {
		t.Fatalf("Game split failed: %v", err)
	}

	hands := player.Hands()
	if len(hands) != 2 {
		t.Errorf("Expected 2 hands after game split, got %d", len(hands))
	}

	// Each hand should have 2 cards after split (original + dealt card)
	for i, hand := range hands {
		if hand.Count() != 2 {
			t.Errorf("Hand %d should have 2 cards after split, got %d", i, hand.Count())
		}
	}
}

func TestSplitBetting(t *testing.T) {
	game := New(1)
	game.AddPlayer("TestPlayer", 1000)
	player := game.GetPlayer("TestPlayer")

	game.StartNewRound()
	player.PlaceBet(100)

	// Set up split scenario
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Seven})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Seven})

	game.PlayerSplit("TestPlayer")

	// Simulate game results for split hands
	game.PayoutResults()

	// The exact result depends on what cards were dealt and dealer's hand
	// But we can verify the betting structure is correct
	if player.Bet() != 0 {
		t.Error("Bet should be cleared after payout")
	}
}

func TestSplitLimitations(t *testing.T) {
	hand := NewHand()

	// Can't split with different ranks
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.King})
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Queen})
	if hand.CanSplit() {
		t.Error("Should not be able to split different ranks")
	}

	// Can't split with one card
	hand.Clear()
	hand.AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	if hand.CanSplit() {
		t.Error("Should not be able to split with one card")
	}

	// Can't split with three cards
	hand.AddCard(cards.Card{Suit: cards.Spades, Rank: cards.Ace})
	hand.AddCard(cards.Card{Suit: cards.Clubs, Rank: cards.Five})
	if hand.CanSplit() {
		t.Error("Should not be able to split with three cards")
	}
}

func TestPlayerSplitInsufficientChips(t *testing.T) {
	player := NewPlayer("TestPlayer", 100)
	player.PlaceBet(100) // All chips

	// Add pair of aces
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Ace})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Ace})

	// Should not be able to split due to insufficient chips
	if player.CanSplit() {
		t.Error("Should not be able to split with insufficient chips")
	}

	err := player.Split()
	if err == nil {
		t.Error("Split should fail with insufficient chips")
	}
}

func TestSplitExample(t *testing.T) {
	player := NewPlayer("Alice", 1000)
	player.PlaceBet(50)

	// Deal a pair of eights
	player.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Eight})
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Eight})

	fmt.Printf("Before split: %d hands\n", len(player.Hands()))
	fmt.Printf("Can split: %t\n", player.CanSplit())

	player.Split()

	fmt.Printf("After split: %d hands\n", len(player.Hands()))
	fmt.Printf("Chips after split: %d\n", player.Chips())
}
