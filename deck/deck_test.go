package deck

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Create new standard deck", func(t *testing.T) {
		deck := New()

		// Standard deck should have 52 cards
		expectedCount := 52
		if len(deck.Cards.Cards) != expectedCount {
			t.Errorf("New deck has %d cards, want %d", len(deck.Cards.Cards), expectedCount)
		}

		// Verify all suits are present
		suitCounts := make(map[Suit]int)
		for _, card := range deck.Cards.Cards {
			suitCounts[card.Suit]++
		}

		expectedSuits := []Suit{Clubs, Diamonds, Hearts, Spades}
		for _, suit := range expectedSuits {
			if suitCounts[suit] != 13 {
				t.Errorf("Suit %s appears %d times, want 13", suit, suitCounts[suit])
			}
		}

		// Verify all ranks are present
		rankCounts := make(map[Rank]int)
		for _, card := range deck.Cards.Cards {
			rankCounts[card.Rank]++
		}

		expectedRanks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
		for _, rank := range expectedRanks {
			if rankCounts[rank] != 4 {
				t.Errorf("Rank %s appears %d times, want 4", rank, rankCounts[rank])
			}
		}
	})

	t.Run("New deck cards are in order", func(t *testing.T) {
		deck := New()

		// Check that cards are in the expected order (suits first, then ranks within each suit)
		suits := []Suit{Clubs, Diamonds, Hearts, Spades}
		ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}

		expectedIndex := 0
		for _, suit := range suits {
			for _, rank := range ranks {
				expectedCard := Card{Rank: rank, Suit: suit}
				actualCard := deck.Cards.Cards[expectedIndex]

				if !reflect.DeepEqual(actualCard, expectedCard) {
					t.Errorf("Card at index %d = %v, want %v", expectedIndex, actualCard, expectedCard)
				}
				expectedIndex++
			}
		}
	})
}

func TestDeck_Inheritance(t *testing.T) {
	t.Run("Deck inherits Cards methods", func(t *testing.T) {
		deck := New()

		// Test that Deck can use Cards methods
		originalCount := len(deck.Cards.Cards)

		// Test Deal method
		card, ok := deck.Deal()
		if !ok {
			t.Errorf("Deal() returned false for new deck")
		}

		expectedFirstCard := Card{Rank: Two, Suit: Clubs}
		if !reflect.DeepEqual(card, expectedFirstCard) {
			t.Errorf("First card dealt = %v, want %v", card, expectedFirstCard)
		}

		if len(deck.Cards.Cards) != originalCount-1 {
			t.Errorf("Deck size after deal = %d, want %d", len(deck.Cards.Cards), originalCount-1)
		}

		// Test Shuffle method
		deck.Shuffle()

		// After shuffle, deck should still have the same number of cards
		if len(deck.Cards.Cards) != originalCount-1 {
			t.Errorf("Deck size after shuffle = %d, want %d", len(deck.Cards.Cards), originalCount-1)
		}
	})
}

func TestDeck_NewDeckStructure(t *testing.T) {
	t.Run("Deck embeds Cards properly", func(t *testing.T) {
		deck := New()

		// Verify that Deck.Cards is accessible
		if deck.Cards.Cards == nil {
			t.Error("Deck.Cards.Cards should not be nil")
		}

		// Verify deck has correct number of cards
		if len(deck.Cards.Cards) != 52 {
			t.Errorf("New deck should have 52 cards, got %d", len(deck.Cards.Cards))
		}
	})
}

func TestDeck_CardOrder(t *testing.T) {
	t.Run("Cards are ordered by suit then rank", func(t *testing.T) {
		deck := New()

		// Verify the order follows: Clubs, Diamonds, Hearts, Spades
		// Within each suit: 2, 3, 4, 5, 6, 7, 8, 9, 10, Jack, Queen, King, Ace

		expectedSuitOrder := []Suit{Clubs, Diamonds, Hearts, Spades}
		expectedRankOrder := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}

		cardIndex := 0
		for _, expectedSuit := range expectedSuitOrder {
			for _, expectedRank := range expectedRankOrder {
				if cardIndex >= len(deck.Cards.Cards) {
					t.Fatalf("Not enough cards in deck at index %d", cardIndex)
				}

				actualCard := deck.Cards.Cards[cardIndex]
				if actualCard.Suit != expectedSuit {
					t.Errorf("Card at index %d has suit %s, expected %s", cardIndex, actualCard.Suit, expectedSuit)
				}
				if actualCard.Rank != expectedRank {
					t.Errorf("Card at index %d has rank %s, expected %s", cardIndex, actualCard.Rank, expectedRank)
				}
				cardIndex++
			}
		}
	})
}

func TestDeck_MethodInheritance(t *testing.T) {
	t.Run("Deck inherits all Cards methods", func(t *testing.T) {
		deck := New()

		// Test Deal method inheritance
		originalLength := len(deck.Cards.Cards)
		card, ok := deck.Deal()
		if !ok {
			t.Error("Deal should return true for non-empty deck")
		}
		if len(deck.Cards.Cards) != originalLength-1 {
			t.Errorf("Deal should reduce deck size by 1, got %d, expected %d", len(deck.Cards.Cards), originalLength-1)
		}

		// Verify the first card dealt is Two of Clubs
		expectedFirstCard := Card{Rank: Two, Suit: Clubs}
		if card != expectedFirstCard {
			t.Errorf("First card should be %v, got %v", expectedFirstCard, card)
		}

		// Test Shuffle method inheritance
		// Store current state
		cardsBeforeShuffle := make([]Card, len(deck.Cards.Cards))
		copy(cardsBeforeShuffle, deck.Cards.Cards)

		deck.Shuffle()

		// Verify same number of cards after shuffle
		if len(deck.Cards.Cards) != len(cardsBeforeShuffle) {
			t.Errorf("Shuffle should not change number of cards, got %d, expected %d", len(deck.Cards.Cards), len(cardsBeforeShuffle))
		}

		// Verify all original cards are still present (though order may have changed)
		cardCounts := make(map[Card]int)
		for _, card := range deck.Cards.Cards {
			cardCounts[card]++
		}

		for _, originalCard := range cardsBeforeShuffle {
			if cardCounts[originalCard] != 1 {
				t.Errorf("Card %v should appear exactly once after shuffle, got %d", originalCard, cardCounts[originalCard])
			}
		}
	})
}

func TestDeck_MultipleDecks(t *testing.T) {
	t.Run("Multiple decks are independent", func(t *testing.T) {
		deck1 := New()
		deck2 := New()

		// Verify both decks have same initial state
		if len(deck1.Cards.Cards) != len(deck2.Cards.Cards) {
			t.Error("Multiple decks should have same initial size")
		}

		// Modify one deck
		card1, _ := deck1.Deal()
		deck2.Shuffle()

		// Verify decks are independent
		if len(deck1.Cards.Cards) == len(deck2.Cards.Cards) {
			t.Error("Dealing from one deck should not affect the other")
		}

		// Verify first deck has the expected card dealt
		expectedFirstCard := Card{Rank: Two, Suit: Clubs}
		if card1 != expectedFirstCard {
			t.Errorf("First card from deck1 should be %v, got %v", expectedFirstCard, card1)
		}

		// Verify second deck still has 52 cards
		if len(deck2.Cards.Cards) != 52 {
			t.Errorf("deck2 should still have 52 cards after shuffle, got %d", len(deck2.Cards.Cards))
		}
	})
}

func TestDeck_EdgeCases(t *testing.T) {
	t.Run("Deal all cards from deck", func(t *testing.T) {
		deck := New()
		dealtCards := []Card{}

		// Deal all 52 cards
		for i := 0; i < 52; i++ {
			card, ok := deck.Deal()
			if !ok {
				t.Errorf("Deal should succeed for card %d", i+1)
			}
			dealtCards = append(dealtCards, card)
		}

		// Verify deck is empty
		if len(deck.Cards.Cards) != 0 {
			t.Errorf("Deck should be empty after dealing all cards, got %d cards remaining", len(deck.Cards.Cards))
		}

		// Try to deal from empty deck
		_, ok := deck.Deal()
		if ok {
			t.Error("Deal should return false for empty deck")
		}

		// Verify we dealt exactly 52 unique cards
		if len(dealtCards) != 52 {
			t.Errorf("Should have dealt 52 cards, got %d", len(dealtCards))
		}

		// Verify all cards are unique
		cardSet := make(map[Card]bool)
		for _, card := range dealtCards {
			if cardSet[card] {
				t.Errorf("Card %v was dealt more than once", card)
			}
			cardSet[card] = true
		}
	})

	t.Run("Shuffle empty deck after dealing all cards", func(t *testing.T) {
		deck := New()

		// Deal all cards
		for len(deck.Cards.Cards) > 0 {
			deck.Deal()
		}

		// Shuffle empty deck (should not panic)
		deck.Shuffle()

		// Verify deck is still empty
		if len(deck.Cards.Cards) != 0 {
			t.Errorf("Empty deck should remain empty after shuffle, got %d cards", len(deck.Cards.Cards))
		}
	})
}
