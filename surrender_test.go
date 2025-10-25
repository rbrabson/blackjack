package blackjack

import (
	"testing"

	"github.com/rbrabson/cards"
)

// TestPlayerSurrender tests basic surrender functionality
func TestPlayerSurrender(t *testing.T) {
	player := NewPlayer("TestPlayer", WithChips(1000))

	// Set up a hand with two cards
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.ClearHand()
	player.hands[0].AddCard(card1)
	player.hands[0].AddCard(card2)

	// Place a bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	chipsAfterBet := player.Chips()

	// Should be able to surrender with 2 cards
	if !player.CanSurrender() {
		t.Fatal("Player should be able to surrender with 2 cards")
	}

	// Surrender the hand
	player.Surrender()

	// Should get half the bet back
	expectedChips := chipsAfterBet + 50 // Half of 100 bet
	if player.Chips() != expectedChips {
		t.Errorf("Expected %d chips after surrender, got %d", expectedChips, player.Chips())
	}

	// Bet should be cleared
	if player.Bet() != 0 {
		t.Errorf("Expected bet to be 0 after surrender, got %d", player.Bet())
	}

	// Hand should be stood
	if !player.CurrentHand().IsStood() {
		t.Error("Hand should be stood after surrender")
	}
}

// TestPlayerCannotSurrenderAfterHit tests that surrender is not allowed after hitting
func TestPlayerCannotSurrenderAfterHit(t *testing.T) {
	player := NewPlayer("TestPlayer", WithChips(1000))

	// Set up a hand with two cards
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Two}

	player.ClearHand()
	player.hands[0].AddCard(card1)
	player.hands[0].AddCard(card2)

	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Should be able to surrender initially
	if !player.CanSurrender() {
		t.Fatal("Player should be able to surrender with 2 cards")
	}

	// Hit (add third card)
	player.Hit(card3)

	// Should no longer be able to surrender
	if player.CanSurrender() {
		t.Error("Player should not be able to surrender after hitting")
	}
}

// TestPlayerCannotSurrenderAfterStand tests that surrender is not allowed after standing
func TestPlayerCannotSurrenderAfterStand(t *testing.T) {
	player := NewPlayer("TestPlayer", WithChips(1000))

	// Set up a hand with two cards
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Nine}

	player.ClearHand()
	player.hands[0].AddCard(card1)
	player.hands[0].AddCard(card2)

	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Should be able to surrender initially
	if !player.CanSurrender() {
		t.Fatal("Player should be able to surrender with 2 cards")
	}

	// Stand
	player.CurrentHand().Stand()

	// Should no longer be able to surrender
	if player.CanSurrender() {
		t.Error("Player should not be able to surrender after standing")
	}
}

// TestPlayerCannotSurrenderWhenBusted tests that surrender is not allowed when busted
func TestPlayerCannotSurrenderWhenBusted(t *testing.T) {
	player := NewPlayer("TestPlayer", WithChips(1000))

	// Set up a busted hand
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Ten}
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Five}

	player.ClearHand()
	player.hands[0].AddCard(card1)
	player.hands[0].AddCard(card2)
	player.hands[0].AddCard(card3) // Busted with 25

	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Should not be able to surrender when busted
	if player.CanSurrender() {
		t.Error("Player should not be able to surrender when busted")
	}
}

// TestGamePlayerSurrender tests surrender through the game interface
func TestGamePlayerSurrender(t *testing.T) {
	game := New(1)
	game.AddPlayer("Alice", WithChips(1000))
	alice := game.GetPlayer("Alice")

	// Set up Alice with a surrenderable hand
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Seven}

	alice.ClearHand()
	alice.hands[0].AddCard(card1)
	alice.hands[0].AddCard(card2)

	err := alice.PlaceBet(200)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	chipsBeforeSurrender := alice.Chips()

	// Surrender through game interface
	err = game.PlayerSurrender("Alice")
	if err != nil {
		t.Fatalf("Surrender failed: %v", err)
	}

	// Should get half bet back
	expectedChips := chipsBeforeSurrender + 100 // Half of 200 bet
	if alice.Chips() != expectedChips {
		t.Errorf("Expected %d chips after surrender, got %d", expectedChips, alice.Chips())
	}

	// Hand should be stood
	if !alice.CurrentHand().IsStood() {
		t.Error("Hand should be stood after surrender")
	}
}

// TestGamePlayerSurrenderInvalidPlayer tests error handling for invalid player
func TestGamePlayerSurrenderInvalidPlayer(t *testing.T) {
	game := New(1)

	err := game.PlayerSurrender("NonexistentPlayer")
	if err == nil {
		t.Error("Expected error for nonexistent player")
	}
}

// TestGamePlayerSurrenderWhenCannotSurrender tests error handling when surrender is not allowed
func TestGamePlayerSurrenderWhenCannotSurrender(t *testing.T) {
	game := New(1)
	game.AddPlayer("Bob", WithChips(1000))
	bob := game.GetPlayer("Bob")

	// Set up Bob with a hand that can't surrender (3 cards)
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Five}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Two}

	bob.ClearHand()
	bob.hands[0].AddCard(card1)
	bob.hands[0].AddCard(card2)
	bob.hands[0].AddCard(card3)

	err := bob.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Should not be able to surrender
	err = game.PlayerSurrender("Bob")
	if err == nil {
		t.Error("Expected error when trying to surrender with 3 cards")
	}
}

// TestSurrenderWithMultipleHands tests surrender behavior with multiple hands (splits)
func TestSurrenderWithMultipleHands(t *testing.T) {
	player := NewPlayer("TestPlayer", WithChips(1000))

	// Set up multiple hands (simulate after a split)
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Eight}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Eight}
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Three}
	card4 := cards.Card{Suit: cards.Diamonds, Rank: cards.Five}

	// Clear and set up hands manually (simulating post-split state)
	player.ClearHand()
	player.hands = append(player.hands, NewHand()) // Add second hand

	// First hand: 8, 3 (should be able to surrender)
	player.hands[0].AddCard(card1)
	player.hands[0].AddCard(card3)

	// Second hand: 8, 5 (should be able to surrender)
	player.hands[1].AddCard(card2)
	player.hands[1].AddCard(card4)

	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}

	// Set current hand to first hand
	player.SetCurrentHandIndex(0)

	// Should not be able to surrender current hand if there are multiple hands
	if player.CanSurrender() {
		t.Fatalf("Player should be able to surrender current hand, numHands=%d", len(player.hands))
	}

}
