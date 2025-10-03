package blackjack

import (
	"strings"
	"testing"

	"github.com/rbrabson/cards"
)

// TestActionTracking tests basic action tracking functionality
func TestActionTracking(t *testing.T) {
	hand := NewHand()

	// Initially no actions
	actions := hand.Actions()
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions initially, got %d", len(actions))
	}

	// Add cards with different action types
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	hand.AddCardWithAction(card1, ActionDeal, "first card")
	hand.AddCardWithAction(card2, ActionDeal, "second card")

	actions = hand.Actions()
	if len(actions) != 2 {
		t.Errorf("Expected 2 actions after dealing, got %d", len(actions))
	}

	// Check first action
	if actions[0].Type != ActionDeal {
		t.Errorf("Expected first action to be deal, got %s", actions[0].Type)
	}
	if actions[0].Card == nil || *actions[0].Card != card1 {
		t.Error("First action should have card1")
	}
	if actions[0].Details != "first card" {
		t.Errorf("Expected details 'first card', got '%s'", actions[0].Details)
	}

	// Record a stand action
	hand.Stand()

	actions = hand.Actions()
	if len(actions) != 3 {
		t.Errorf("Expected 3 actions after standing, got %d", len(actions))
	}

	if actions[2].Type != ActionStand {
		t.Errorf("Expected third action to be stand, got %s", actions[2].Type)
	}
	if actions[2].Card != nil {
		t.Error("Stand action should not have a card")
	}
}

// TestPlayerActionTracking tests action tracking through player methods
func TestPlayerActionTracking(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Deal initial cards
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.DealCard(card1)
	player.DealCard(card2)

	// Check initial deal actions
	actions := player.CurrentHand().Actions()
	if len(actions) != 2 {
		t.Errorf("Expected 2 deal actions, got %d", len(actions))
	}

	if actions[0].Type != ActionDeal || actions[1].Type != ActionDeal {
		t.Error("Both initial actions should be deals")
	}

	// Test hit action
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Five}
	player.Hit(card3)

	actions = player.CurrentHand().Actions()
	if len(actions) != 3 {
		t.Errorf("Expected 3 actions after hit, got %d", len(actions))
	}

	if actions[2].Type != ActionHit {
		t.Errorf("Expected third action to be hit, got %s", actions[2].Type)
	}

	// Test stand action
	player.CurrentHand().Stand()

	actions = player.CurrentHand().Actions()
	if len(actions) != 4 {
		t.Errorf("Expected 4 actions after stand, got %d", len(actions))
	}

	if actions[3].Type != ActionStand {
		t.Errorf("Expected fourth action to be stand, got %s", actions[3].Type)
	}
}

// TestSurrenderActionTracking tests surrender action tracking
func TestSurrenderActionTracking(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up a hand for surrender
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.DealCard(card1)
	player.DealCard(card2)
	player.PlaceBet(100)

	// Surrender
	player.Surrender()

	actions := player.CurrentHand().Actions()
	if len(actions) != 4 { // 2 deals + 1 surrender + 1 stand
		t.Errorf("Expected 4 actions after surrender, got %d", len(actions))
	}

	// Find surrender action
	surrenderFound := false
	standFound := false
	for _, action := range actions {
		if action.Type == ActionSurrender {
			surrenderFound = true
			if !strings.Contains(action.Details, "received") && !strings.Contains(action.Details, "chips") {
				t.Error("Surrender action should include chip return details")
			}
		}
		if action.Type == ActionStand {
			standFound = true
		}
	}

	if !surrenderFound {
		t.Error("Surrender action not found")
	}
	if !standFound {
		t.Error("Stand action should be recorded after surrender")
	}
}

// TestDoubleDownActionTracking tests double down action tracking
func TestDoubleDownActionTracking(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up for double down
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	player.DealCard(card1)
	player.DealCard(card2)
	player.PlaceBet(100)

	// Double down
	err := player.DoubleDown()
	if err != nil {
		t.Fatalf("DoubleDown failed: %v", err)
	}

	// Add double down card
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Five}
	player.DoubleDownHit(card3)

	actions := player.CurrentHand().Actions()
	if len(actions) != 4 { // 2 deals + 1 double + 1 double hit
		t.Errorf("Expected 4 actions after double down, got %d", len(actions))
	}

	// Check for double action
	doubleFound := false
	doubleHitFound := false
	for _, action := range actions {
		if action.Type == ActionDouble && action.Card == nil {
			doubleFound = true
			if !strings.Contains(action.Details, "bet increased") {
				t.Error("Double action should include bet increase details")
			}
		}
		if action.Type == ActionDouble && action.Card != nil {
			doubleHitFound = true
			if *action.Card != card3 {
				t.Error("Double hit action should have the correct card")
			}
		}
	}

	if !doubleFound {
		t.Error("Double action not found")
	}
	if !doubleHitFound {
		t.Error("Double hit action not found")
	}
}

// TestSplitActionTracking tests split action tracking
func TestSplitActionTracking(t *testing.T) {
	player := NewPlayer("TestPlayer", 1000)

	// Set up for split
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Eight}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Eight}

	player.DealCard(card1)
	player.DealCard(card2)
	player.PlaceBet(100)

	// Split
	err := player.Split()
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Check first hand actions
	actions1 := player.hands[0].Actions()
	splitFound := false
	for _, action := range actions1 {
		if action.Type == ActionSplit {
			splitFound = true
			break
		}
	}
	if !splitFound {
		t.Error("Split action not found in first hand")
	}

	// Check second hand actions
	actions2 := player.hands[1].Actions()
	splitFoundIn2 := false
	for _, action := range actions2 {
		if action.Type == ActionSplit {
			splitFoundIn2 = true
			if !strings.Contains(action.Details, "created from split") {
				t.Error("Second hand should have 'created from split' details")
			}
			break
		}
	}
	if !splitFoundIn2 {
		t.Error("Split action not found in second hand")
	}
}

// TestDealerActionTracking tests dealer action tracking
func TestDealerActionTracking(t *testing.T) {
	dealer := NewDealer()

	// Deal initial cards
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	dealer.DealCard(card1)
	dealer.DealCard(card2)

	actions := dealer.Hand().Actions()
	if len(actions) != 2 {
		t.Errorf("Expected 2 deal actions, got %d", len(actions))
	}

	if actions[0].Type != ActionDeal || actions[1].Type != ActionDeal {
		t.Error("Both initial actions should be deals")
	}

	// Dealer hits
	card3 := cards.Card{Suit: cards.Clubs, Rank: cards.Five}
	dealer.Hit(card3)

	actions = dealer.Hand().Actions()
	if len(actions) != 3 {
		t.Errorf("Expected 3 actions after hit, got %d", len(actions))
	}

	if actions[2].Type != ActionHit {
		t.Errorf("Expected third action to be hit, got %s", actions[2].Type)
	}

	// Dealer stands
	dealer.Stand()

	actions = dealer.Hand().Actions()
	if len(actions) != 4 {
		t.Errorf("Expected 4 actions after stand, got %d", len(actions))
	}

	if actions[3].Type != ActionStand {
		t.Errorf("Expected fourth action to be stand, got %s", actions[3].Type)
	}
}

// TestActionSummary tests the action summary string generation
func TestActionSummary(t *testing.T) {
	hand := NewHand()

	// Test empty hand
	summary := hand.ActionSummary()
	if summary != "No actions" {
		t.Errorf("Expected 'No actions' for empty hand, got '%s'", summary)
	}

	// Add some actions
	card1 := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	card2 := cards.Card{Suit: cards.Hearts, Rank: cards.Six}

	hand.AddCardWithAction(card1, ActionDeal, "initial")
	hand.AddCardWithAction(card2, ActionDeal, "initial")
	hand.AddCardWithAction(cards.Card{Suit: cards.Clubs, Rank: cards.Five}, ActionHit, "")
	hand.RecordAction(ActionStand, "")

	summary = hand.ActionSummary()
	expectedParts := []string{"dealt Ten of Spades", "dealt Six of Hearts", "hit Five of Clubs", "stand"}

	for _, part := range expectedParts {
		if !strings.Contains(summary, part) {
			t.Errorf("Expected summary to contain '%s', got '%s'", part, summary)
		}
	}
}

// TestActionTimestamps tests that actions have timestamps
func TestActionTimestamps(t *testing.T) {
	hand := NewHand()

	card := cards.Card{Suit: cards.Spades, Rank: cards.Ten}
	hand.AddCardWithAction(card, ActionDeal, "test")

	actions := hand.Actions()
	if len(actions) != 1 {
		t.Fatalf("Expected 1 action, got %d", len(actions))
	}

	if actions[0].Timestamp.IsZero() {
		t.Error("Action timestamp should not be zero")
	}
}
