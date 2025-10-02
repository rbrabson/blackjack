package blackjack

import (
	"testing"

	"github.com/rbrabson/cards"
)

// TestFourSplitScenario tests a realistic scenario where a player splits 4 times
func TestFourSplitScenario(t *testing.T) {
	// Create a game and add a player with plenty of chips
	game := New(1)
	game.AddPlayer("Alice", 10000)
	alice := game.GetPlayer("Alice")
	
	// Set up Alice with four aces (extremely rare but possible scenario)
	ace1 := cards.Card{Suit: cards.Spades, Rank: cards.Ace}
	ace2 := cards.Card{Suit: cards.Hearts, Rank: cards.Ace}
	ace3 := cards.Card{Suit: cards.Clubs, Rank: cards.Ace}
	ace4 := cards.Card{Suit: cards.Diamonds, Rank: cards.Ace}
	
	alice.ClearHand()
	alice.hands[0].AddCard(ace1)
	alice.hands[0].AddCard(ace2)
	
	// Place initial bet
	err := alice.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}
	
	chipsAfterBet := alice.Chips()
	
	// First split: AA -> A|A
	if !alice.CanSplit() {
		t.Fatal("Alice should be able to split her aces")
	}
	
	err = alice.Split()
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}
	
	if alice.Chips() != chipsAfterBet-100 {
		t.Errorf("Expected %d chips after first split, got %d", chipsAfterBet-100, alice.Chips())
	}
	
	if len(alice.hands) != 2 {
		t.Errorf("Expected 2 hands after first split, got %d", len(alice.hands))
	}
	
	// Set up both hands with pairs for further splitting
	alice.hands[0] = *NewHand()
	alice.hands[0].AddCard(ace1)
	alice.hands[0].AddCard(ace3)
	
	alice.hands[1] = *NewHand()
	alice.hands[1].AddCard(ace2)
	alice.hands[1].AddCard(ace4)
	
	// Second split: split the first hand
	alice.SetCurrentHandIndex(0)
	if !alice.CanSplit() {
		t.Fatal("Alice should be able to split again")
	}
	
	err = alice.Split()
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}
	
	if len(alice.hands) != 3 {
		t.Errorf("Expected 3 hands after second split, got %d", len(alice.hands))
	}
	
	// Set up hands for third split
	ace5 := cards.Card{Suit: cards.Spades, Rank: cards.Ace}
	ace6 := cards.Card{Suit: cards.Hearts, Rank: cards.Ace}
	
	alice.hands[1] = *NewHand()
	alice.hands[1].AddCard(ace2)
	alice.hands[1].AddCard(ace5)
	
	// Third split: split the second hand to get 4 hands total
	alice.SetCurrentHandIndex(1)
	if !alice.CanSplit() {
		t.Fatal("Alice should be able to split a third time")
	}
	
	err = alice.Split()
	if err != nil {
		t.Fatalf("Third split failed: %v", err)
	}
	
	if len(alice.hands) != 4 {
		t.Errorf("Expected 4 hands after third split, got %d", len(alice.hands))
	}
	
	// Set up one hand to be splittable for testing the limit
	alice.hands[2] = *NewHand()
	alice.hands[2].AddCard(ace3)
	alice.hands[2].AddCard(ace6)
	
	// Fourth split attempt: should fail because we already have 4 hands
	alice.SetCurrentHandIndex(2)
	if alice.CanSplit() {
		t.Error("Alice should NOT be able to split when she already has 4 hands")
	}
	
	// Attempt to split should fail
	err = alice.Split()
	if err == nil {
		t.Error("Split should have failed - already at 4 hand limit")
	}
	
	// Should still have exactly 4 hands
	if len(alice.hands) != 4 {
		t.Errorf("Should still have exactly 4 hands, got %d", len(alice.hands))
	}
	
	// Verify we spent the right amount (initial bet + 3 split bets)
	expectedChips := 10000 - 400 // 100 for initial bet + 100 each for 3 splits
	if alice.Chips() != expectedChips {
		t.Errorf("Expected %d chips, got %d", expectedChips, alice.Chips())
	}
}

// TestSplitLimitInGame tests the split limit within a game context
func TestSplitLimitInGame(t *testing.T) {
	game := New(1)
	game.AddPlayer("Bob", 5000)
	bob := game.GetPlayer("Bob")
	
	// Clear Bob's hand and give him a splittable pair
	bob.ClearHand()
	king1 := cards.Card{Suit: cards.Spades, Rank: cards.King}
	king2 := cards.Card{Suit: cards.Hearts, Rank: cards.King}
	bob.hands[0].AddCard(king1)
	bob.hands[0].AddCard(king2)
	
	err := bob.PlaceBet(50)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}
	
	// Test splits through the game interface
	// First split
	err = game.PlayerSplit("Bob")
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}
	
	if len(bob.hands) != 2 {
		t.Errorf("Expected 2 hands after first split, got %d", len(bob.hands))
	}
	
	// Set up hands for more splits
	bob.hands[0] = *NewHand()
	bob.hands[0].AddCard(king1)
	bob.hands[0].AddCard(cards.Card{Suit: cards.Clubs, Rank: cards.King})
	
	bob.hands[1] = *NewHand()
	bob.hands[1].AddCard(king2)
	bob.hands[1].AddCard(cards.Card{Suit: cards.Diamonds, Rank: cards.King})
	
	// Second split
	bob.SetCurrentHandIndex(0)
	err = game.PlayerSplit("Bob")
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}
	
	if len(bob.hands) != 3 {
		t.Errorf("Expected 3 hands after second split, got %d", len(bob.hands))
	}
	
	// Set up for third split
	bob.hands[1] = *NewHand()
	bob.hands[1].AddCard(cards.Card{Suit: cards.Spades, Rank: cards.King})
	bob.hands[1].AddCard(cards.Card{Suit: cards.Hearts, Rank: cards.King})
	
	// Third split
	bob.SetCurrentHandIndex(1)
	err = game.PlayerSplit("Bob")
	if err != nil {
		t.Fatalf("Third split failed: %v", err)
	}
	
	if len(bob.hands) != 4 {
		t.Errorf("Expected 4 hands after third split, got %d", len(bob.hands))
	}
	
	// Set up one more hand to try splitting past the limit
	bob.hands[2] = *NewHand()
	bob.hands[2].AddCard(cards.Card{Suit: cards.Clubs, Rank: cards.King})
	bob.hands[2].AddCard(cards.Card{Suit: cards.Diamonds, Rank: cards.King})
	
	// Fourth split attempt should fail
	bob.SetCurrentHandIndex(2)
	err = game.PlayerSplit("Bob")
	if err == nil {
		t.Error("Fourth split should have failed due to 4-hand limit")
	}
	
	// Should still have exactly 4 hands
	if len(bob.hands) != 4 {
		t.Errorf("Should still have exactly 4 hands after failed split, got %d", len(bob.hands))
	}
}