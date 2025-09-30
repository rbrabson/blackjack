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
