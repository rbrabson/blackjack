package blackjack

import (
	"testing"

	"github.com/rbrabson/cards"
)

// TestPlayerSplitLimit tests that a player can split up to 4 hands but no more
func TestPlayerSplitLimit(t *testing.T) {
	player := NewPlayer("TestPlayer", 10000) // Give plenty of chips
	
	// Set up the first hand with a pair that can be split
	ace1 := cards.Card{Suit: cards.Spades, Rank: cards.Ace}
	ace2 := cards.Card{Suit: cards.Hearts, Rank: cards.Ace}
	
	player.ClearHand()
	player.hands[0].AddCard(ace1)
	player.hands[0].AddCard(ace2)
	
	// Place a bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}
	
	// First split: 1 hand -> 2 hands
	if !player.CanSplit() {
		t.Fatal("Player should be able to split first time")
	}
	
	err = player.Split()
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}
	
	if len(player.hands) != 2 {
		t.Errorf("Expected 2 hands after first split, got %d", len(player.hands))
	}
	
	// Add another ace to each hand to make them splittable again
	ace3 := cards.Card{Suit: cards.Clubs, Rank: cards.Ace}
	ace4 := cards.Card{Suit: cards.Diamonds, Rank: cards.Ace}
	
	player.hands[0].AddCard(ace3)
	player.hands[1].AddCard(ace4)
	
	// Set current hand to first hand for second split
	player.SetCurrentHandIndex(0)
	
	// Second split: 2 hands -> 3 hands
	if !player.CanSplit() {
		t.Fatal("Player should be able to split second time")
	}
	
	err = player.Split()
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}
	
	if len(player.hands) != 3 {
		t.Errorf("Expected 3 hands after second split, got %d", len(player.hands))
	}
	
	// Add ace to the new third hand to make it splittable
	ace7 := cards.Card{Suit: cards.Clubs, Rank: cards.Ace}
	player.hands[2].AddCard(ace7)
	
	// Set current hand to second hand for third split
	player.SetCurrentHandIndex(1)
	
	// Third split: 3 hands -> 4 hands
	if !player.CanSplit() {
		t.Fatal("Player should be able to split third time")
	}
	
	err = player.Split()
	if err != nil {
		t.Fatalf("Third split failed: %v", err)
	}
	
	if len(player.hands) != 4 {
		t.Errorf("Expected 4 hands after third split, got %d", len(player.hands))
	}
	
	// Add ace to the new fourth hand
	ace8 := cards.Card{Suit: cards.Diamonds, Rank: cards.Ace}
	player.hands[3].AddCard(ace8)
	
	// Try fourth split: should fail because we already have 4 hands
	player.SetCurrentHandIndex(2) // Try to split the third hand
	
	if player.CanSplit() {
		t.Error("Player should NOT be able to split when already having 4 hands")
	}
	
	err = player.Split()
	if err == nil {
		t.Error("Split should have failed when player already has 4 hands")
	}
	
	if len(player.hands) != 4 {
		t.Errorf("Should still have exactly 4 hands, got %d", len(player.hands))
	}
}

// TestPlayerSplitLimitWithInsufficientChips tests the split limit interaction with chip constraints
func TestPlayerSplitLimitWithInsufficientChips(t *testing.T) {
	player := NewPlayer("TestPlayer", 300) // Limited chips for only 2 additional splits
	
	// Set up hand with a pair
	king1 := cards.Card{Suit: cards.Spades, Rank: cards.King}
	king2 := cards.Card{Suit: cards.Hearts, Rank: cards.King}
	
	player.ClearHand()
	player.hands[0].AddCard(king1)
	player.hands[0].AddCard(king2)
	
	// Place a bet
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}
	
	// First split should work (have 200 chips remaining)
	if !player.CanSplit() {
		t.Fatal("Player should be able to split first time")
	}
	
	err = player.Split()
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}
	
	// Add kings to make hands splittable again
	king3 := cards.Card{Suit: cards.Clubs, Rank: cards.King}
	king4 := cards.Card{Suit: cards.Diamonds, Rank: cards.King}
	
	player.hands[0].AddCard(king3)
	player.hands[1].AddCard(king4)
	
	// Second split should work (have 100 chips remaining)
	player.SetCurrentHandIndex(0)
	if !player.CanSplit() {
		t.Fatal("Player should be able to split second time")
	}
	
	err = player.Split()
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}
	
	// Add kings to make hands splittable again
	king5 := cards.Card{Suit: cards.Spades, Rank: cards.King}
	king6 := cards.Card{Suit: cards.Hearts, Rank: cards.King}
	
	player.hands[1].AddCard(king5)
	player.hands[2].AddCard(king6)
	
	// Third split should fail due to insufficient chips (have 0 chips remaining)
	player.SetCurrentHandIndex(1)
	if player.CanSplit() {
		t.Error("Player should NOT be able to split due to insufficient chips")
	}
	
	if len(player.hands) != 3 {
		t.Errorf("Should have exactly 3 hands, got %d", len(player.hands))
	}
}

// TestPlayerSplitLimitEdgeCases tests edge cases around the split limit
func TestPlayerSplitLimitEdgeCases(t *testing.T) {
	player := NewPlayer("TestPlayer", 10000)
	
	// Test that we can't split if current hand is not splittable
	two1 := cards.Card{Suit: cards.Spades, Rank: cards.Two}
	three1 := cards.Card{Suit: cards.Hearts, Rank: cards.Three}
	
	player.ClearHand()
	player.hands[0].AddCard(two1)
	player.hands[0].AddCard(three1)
	
	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Failed to place bet: %v", err)
	}
	
	if player.CanSplit() {
		t.Error("Player should not be able to split non-matching cards")
	}
	
	// Test that we start fresh each round
	player.ClearHand() // This should reset to 1 hand
	if len(player.hands) != 1 {
		t.Errorf("ClearHand should reset to 1 hand, got %d", len(player.hands))
	}
	
	// Now we should be able to split again if we have a pair
	ace1 := cards.Card{Suit: cards.Spades, Rank: cards.Ace}
	ace2 := cards.Card{Suit: cards.Hearts, Rank: cards.Ace}
	
	player.hands[0].AddCard(ace1)
	player.hands[0].AddCard(ace2)
	
	if !player.CanSplit() {
		t.Error("Player should be able to split again after ClearHand")
	}
}