package blackjack

import (
	"fmt"
	"testing"
)

// TestChipManagerInterface demonstrates using a custom chip manager
func TestChipManagerInterface(t *testing.T) {
	// Test with default chip manager
	player1 := NewPlayer("Player1", 1000)
	if player1.Chips() != 1000 {
		t.Errorf("Expected 1000 chips, got %d", player1.Chips())
	}

	// Test with custom chip manager
	customChipManager := NewDefaultChipManager(500)
	player2 := NewPlayerWithChipManager("Player2", customChipManager)
	if player2.Chips() != 500 {
		t.Errorf("Expected 500 chips, got %d", player2.Chips())
	}

	// Test chip operations through the interface
	err := player2.PlaceBet(100)
	if err != nil {
		t.Errorf("Unexpected error placing bet: %v", err)
	}
	if player2.Chips() != 400 {
		t.Errorf("Expected 400 chips after bet, got %d", player2.Chips())
	}

	// Test win bet
	player2.WinBet(1.5)                       // 1.5x multiplier
	expectedChips := 400 + 100 + int(100*1.5) // 400 + 100 (original bet) + 150 (winnings)
	if player2.Chips() != expectedChips {
		t.Errorf("Expected %d chips after win, got %d", expectedChips, player2.Chips())
	}
}

// TrackingChipManager is a custom chip manager that tracks operations
type TrackingChipManager struct {
	chips          int
	operationCount int
}

// GetChips returns the current chip count
func (t *TrackingChipManager) GetChips() int {
	return t.chips
}

// SetChips sets the chip count to the specified amount
func (t *TrackingChipManager) SetChips(amount int) {
	t.operationCount++
	t.chips = amount
}

// AddChips adds the specified amount to the chip count
func (t *TrackingChipManager) AddChips(amount int) {
	t.operationCount++
	t.chips += amount
}

// DeductChips removes the specified amount from the chip count
func (t *TrackingChipManager) DeductChips(amount int) error {
	t.operationCount++
	if amount > t.chips {
		return fmt.Errorf("insufficient chips")
	}
	t.chips -= amount
	return nil
}

// HasEnoughChips returns true if there are enough chips for the specified amount
func (t *TrackingChipManager) HasEnoughChips(amount int) bool {
	return t.chips >= amount
}

// TestCustomChipManager demonstrates creating a custom chip manager implementation
func TestCustomChipManager(t *testing.T) {
	trackingManager := &TrackingChipManager{chips: 1000, operationCount: 0}
	player := NewPlayerWithChipManager("TrackingPlayer", trackingManager)

	// Place a bet (should increment operation count)
	err := player.PlaceBet(100)
	if err != nil {
		t.Errorf("Unexpected error placing bet: %v", err)
	}

	// Win the bet (should increment operation count)
	player.WinBet(1.0)

	if trackingManager.operationCount != 2 {
		t.Errorf("Expected 2 operations tracked, got %d", trackingManager.operationCount)
	}
}
