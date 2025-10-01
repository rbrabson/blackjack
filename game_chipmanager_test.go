package blackjack

import (
	"testing"
)

// TestGameAddPlayerWithChipManager tests adding a player with a custom chip manager
func TestGameAddPlayerWithChipManager(t *testing.T) {
	game := New(1)

	// Add a player with default chip manager
	game.AddPlayer("Alice", 1000)
	alice := game.GetPlayer("Alice")
	if alice == nil {
		t.Fatal("Alice not found in game")
	}
	if alice.Chips() != 1000 {
		t.Errorf("Expected Alice to have 1000 chips, got %d", alice.Chips())
	}

	// Add a player with custom chip manager
	customChipManager := &TrackingChipManager{chips: 500, operationCount: 0}
	game.AddPlayer("Bob", 500, WithChipManager(customChipManager))
	bob := game.GetPlayer("Bob")
	if bob == nil {
		t.Fatal("Bob not found in game")
	}
	if bob.Chips() != 500 {
		t.Errorf("Expected Bob to have 500 chips, got %d", bob.Chips())
	}

	// Test that Bob's chip manager is the custom one by checking operation tracking
	customChipManager.operationCount = 0 // Reset counter after player creation
	err := bob.PlaceBet(100)
	if err != nil {
		t.Errorf("Unexpected error placing bet: %v", err)
	}

	if customChipManager.operationCount != 1 {
		t.Errorf("Expected 1 operation tracked, got %d", customChipManager.operationCount)
	}

	// Verify both players are in the game
	if len(game.Players()) != 2 {
		t.Errorf("Expected 2 players in game, got %d", len(game.Players()))
	}
}

// TestGameAddPlayerBackwardCompatibility ensures the original AddPlayer method still works
func TestGameAddPlayerBackwardCompatibility(t *testing.T) {
	game := New(1)

	// This should work exactly as before
	game.AddPlayer("Charlie", 750)
	charlie := game.GetPlayer("Charlie")
	if charlie == nil {
		t.Fatal("Charlie not found in game")
	}
	if charlie.Chips() != 750 {
		t.Errorf("Expected Charlie to have 750 chips, got %d", charlie.Chips())
	}

	// Should be able to place bets normally
	err := charlie.PlaceBet(50)
	if err != nil {
		t.Errorf("Unexpected error placing bet: %v", err)
	}
	if charlie.Chips() != 700 {
		t.Errorf("Expected Charlie to have 700 chips after bet, got %d", charlie.Chips())
	}
}
