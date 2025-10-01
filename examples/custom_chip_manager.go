package main

import (
	"fmt"

	"github.com/rbrabson/blackjack"
)

// ExampleCustomChipManager demonstrates a chip manager with daily limits
type ExampleCustomChipManager struct {
	chips      int
	dailySpent int
	dailyLimit int
}

func NewExampleCustomChipManager(initialChips, dailyLimit int) *ExampleCustomChipManager {
	return &ExampleCustomChipManager{
		chips:      initialChips,
		dailySpent: 0,
		dailyLimit: dailyLimit,
	}
}

func (e *ExampleCustomChipManager) GetChips() int {
	return e.chips
}

func (e *ExampleCustomChipManager) SetChips(amount int) {
	e.chips = amount
}

func (e *ExampleCustomChipManager) AddChips(amount int) {
	e.chips += amount
}

func (e *ExampleCustomChipManager) DeductChips(amount int) error {
	if amount > e.chips {
		return fmt.Errorf("insufficient chips: have %d, need %d", e.chips, amount)
	}
	if e.dailySpent+amount > e.dailyLimit {
		return fmt.Errorf("daily limit exceeded: spent %d, limit %d, trying to spend %d more",
			e.dailySpent, e.dailyLimit, amount)
	}
	e.chips -= amount
	e.dailySpent += amount
	return nil
}

func (e *ExampleCustomChipManager) HasEnoughChips(amount int) bool {
	return e.chips >= amount && e.dailySpent+amount <= e.dailyLimit
}

func main() {
	// Create a new game
	game := blackjack.New(6)

	// Add a regular player
	game.AddPlayer("Alice", 1000)

	// Add a player with daily spending limits
	limitedChipManager := NewExampleCustomChipManager(1000, 500) // $500 daily limit
	game.AddPlayer("Bob", 1000, blackjack.WithChipManager(limitedChipManager))

	fmt.Println("Game created with 2 players:")
	for _, player := range game.Players() {
		fmt.Printf("- %s: %d chips\n", player.Name(), player.Chips())
	}

	// Demonstrate the daily limit feature
	bob := game.GetPlayer("Bob")

	fmt.Println("\nTesting Bob's daily limit...")

	// This should work (under limit)
	err := bob.PlaceBet(300)
	if err != nil {
		fmt.Printf("Error placing 300 chip bet: %v\n", err)
	} else {
		fmt.Printf("Successfully placed 300 chip bet. Remaining chips: %d\n", bob.Chips())
		bob.LoseBet() // Simulate losing the bet
	}

	// This should fail (exceeds daily limit)
	err = bob.PlaceBet(300)
	if err != nil {
		fmt.Printf("Error placing second 300 chip bet: %v\n", err)
	} else {
		fmt.Printf("Successfully placed second 300 chip bet. Remaining chips: %d\n", bob.Chips())
	}

	// Alice shouldn't have this limitation
	alice := game.GetPlayer("Alice")
	err = alice.PlaceBet(600)
	if err != nil {
		fmt.Printf("Error with Alice's 600 chip bet: %v\n", err)
	} else {
		fmt.Printf("Alice successfully placed 600 chip bet. Remaining chips: %d\n", alice.Chips())
	}
}
