package blackjack

import "fmt"

// ChipManager interface defines the operations for managing player chips
type ChipManager interface {
	GetChips() int                  // GetChips returns the current chip count
	SetChips(amount int)            // SetChips sets the chip count to the specified amount
	AddChips(amount int)            // AddChips adds the specified amount to the chip count
	DeductChips(amount int) error   // DeductChips removes the specified amount from the chip count
	HasEnoughChips(amount int) bool // HasEnoughChips returns true if there are enough chips for the specified amount
}

// DefaultChipManager implements ChipManager with simple integer-based chip management
type DefaultChipManager struct {
	chips int
}

// NewDefaultChipManager creates a new default chip manager with the given initial amount
func NewDefaultChipManager(initialChips int) *DefaultChipManager {
	return &DefaultChipManager{chips: initialChips}
}

// GetChips returns the current chip count
func (c *DefaultChipManager) GetChips() int {
	return c.chips
}

// SetChips sets the chip count to the specified amount
func (c *DefaultChipManager) SetChips(amount int) {
	c.chips = amount
}

// AddChips adds the specified amount to the chip count
func (c *DefaultChipManager) AddChips(amount int) {
	c.chips += amount
}

// DeductChips removes the specified amount from the chip count
func (c *DefaultChipManager) DeductChips(amount int) error {
	if amount > c.chips {
		return fmt.Errorf("insufficient chips: have %d, need %d", c.chips, amount)
	}
	c.chips -= amount
	return nil
}

// HasEnoughChips returns true if there are enough chips for the specified amount
func (c *DefaultChipManager) HasEnoughChips(amount int) bool {
	return c.chips >= amount
}
