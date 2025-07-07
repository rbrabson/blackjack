package deck

import (
	"reflect"
	"testing"
)

func TestCard(t *testing.T) {
	tests := []struct {
		name     string
		rank     Rank
		suit     Suit
		expected Card
	}{
		{
			name:     "Create Ace of Spades",
			rank:     Ace,
			suit:     Spades,
			expected: Card{Rank: Ace, Suit: Spades},
		},
		{
			name:     "Create King of Hearts",
			rank:     King,
			suit:     Hearts,
			expected: Card{Rank: King, Suit: Hearts},
		},
		{
			name:     "Create Two of Clubs",
			rank:     Two,
			suit:     Clubs,
			expected: Card{Rank: Two, Suit: Clubs},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := Card{
				Rank: tt.rank,
				Suit: tt.suit,
			}
			if !reflect.DeepEqual(card, tt.expected) {
				t.Errorf("Card = %v, want %v", card, tt.expected)
			}
		})
	}
}

func TestCards_Deal(t *testing.T) {
	t.Run("Deal from non-empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{
				{Rank: Ace, Suit: Spades},
				{Rank: King, Suit: Hearts},
				{Rank: Queen, Suit: Clubs},
			},
		}

		card, ok := cards.Deal()
		if !ok {
			t.Errorf("Deal() returned false for non-empty deck")
		}

		expectedCard := Card{Rank: Ace, Suit: Spades}
		if !reflect.DeepEqual(card, expectedCard) {
			t.Errorf("Deal() = %v, want %v", card, expectedCard)
		}

		if len(cards.Cards) != 2 {
			t.Errorf("Cards length after deal = %d, want 2", len(cards.Cards))
		}

		// Check that the first card was removed
		if reflect.DeepEqual(cards.Cards[0], expectedCard) {
			t.Errorf("First card was not removed from deck")
		}
	})

	t.Run("Deal from empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{},
		}

		card, ok := cards.Deal()
		if ok {
			t.Errorf("Deal() returned true for empty deck")
		}

		expectedCard := Card{}
		if !reflect.DeepEqual(card, expectedCard) {
			t.Errorf("Deal() from empty deck = %v, want %v", card, expectedCard)
		}
	})

	t.Run("Deal until empty", func(t *testing.T) {
		originalCards := []Card{
			{Rank: Ace, Suit: Spades},
			{Rank: King, Suit: Hearts},
		}
		cards := &Cards{
			Cards: make([]Card, len(originalCards)),
		}
		copy(cards.Cards, originalCards)

		// Deal first card
		card1, ok1 := cards.Deal()
		if !ok1 {
			t.Errorf("First deal returned false")
		}
		if !reflect.DeepEqual(card1, originalCards[0]) {
			t.Errorf("First deal = %v, want %v", card1, originalCards[0])
		}

		// Deal second card
		card2, ok2 := cards.Deal()
		if !ok2 {
			t.Errorf("Second deal returned false")
		}
		if !reflect.DeepEqual(card2, originalCards[1]) {
			t.Errorf("Second deal = %v, want %v", card2, originalCards[1])
		}

		// Try to deal from empty deck
		_, ok3 := cards.Deal()
		if ok3 {
			t.Errorf("Deal from empty deck returned true")
		}
	})
}

func TestCards_Shuffle(t *testing.T) {
	t.Run("Shuffle changes order", func(t *testing.T) {
		// Create a deck with multiple cards
		originalCards := []Card{
			{Rank: Ace, Suit: Spades},
			{Rank: King, Suit: Hearts},
			{Rank: Queen, Suit: Clubs},
			{Rank: Jack, Suit: Diamonds},
			{Rank: Ten, Suit: Spades},
			{Rank: Nine, Suit: Hearts},
			{Rank: Eight, Suit: Clubs},
			{Rank: Seven, Suit: Diamonds},
		}

		cards := &Cards{
			Cards: make([]Card, len(originalCards)),
		}
		copy(cards.Cards, originalCards)

		// Shuffle multiple times to increase likelihood of order change
		var shuffled bool
		for i := 0; i < 10; i++ {
			cards.Shuffle()
			if !reflect.DeepEqual(cards.Cards, originalCards) {
				shuffled = true
				break
			}
			// Reset for next iteration
			copy(cards.Cards, originalCards)
		}

		if !shuffled {
			t.Log("Warning: Shuffle did not change order after 10 attempts (this could happen by chance)")
		}

		// Verify all cards are still present
		if len(cards.Cards) != len(originalCards) {
			t.Errorf("Shuffle changed number of cards: got %d, want %d", len(cards.Cards), len(originalCards))
		}

		// Verify no cards were lost or duplicated
		cardCounts := make(map[Card]int)
		for _, card := range cards.Cards {
			cardCounts[card]++
		}

		for _, originalCard := range originalCards {
			if cardCounts[originalCard] != 1 {
				t.Errorf("Card %v appears %d times after shuffle, want 1", originalCard, cardCounts[originalCard])
			}
		}
	})

	t.Run("Shuffle empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{},
		}

		// Should not panic
		cards.Shuffle()

		if len(cards.Cards) != 0 {
			t.Errorf("Shuffle of empty deck changed length: got %d, want 0", len(cards.Cards))
		}
	})

	t.Run("Shuffle single card", func(t *testing.T) {
		originalCard := Card{Rank: Ace, Suit: Spades}
		cards := &Cards{
			Cards: []Card{originalCard},
		}

		cards.Shuffle()

		if len(cards.Cards) != 1 {
			t.Errorf("Shuffle of single card changed length: got %d, want 1", len(cards.Cards))
		}

		if !reflect.DeepEqual(cards.Cards[0], originalCard) {
			t.Errorf("Shuffle changed single card: got %v, want %v", cards.Cards[0], originalCard)
		}
	})
}

func TestSuitConstants(t *testing.T) {
	tests := []struct {
		name     string
		suit     Suit
		expected string
	}{
		{"Clubs", Clubs, "Clubs"},
		{"Diamonds", Diamonds, "Diamonds"},
		{"Hearts", Hearts, "Hearts"},
		{"Spades", Spades, "Spades"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.suit) != tt.expected {
				t.Errorf("Suit %s = %s, want %s", tt.name, string(tt.suit), tt.expected)
			}
		})
	}
}

func TestRankConstants(t *testing.T) {
	tests := []struct {
		name     string
		rank     Rank
		expected string
	}{
		{"Two", Two, "2"},
		{"Three", Three, "3"},
		{"Four", Four, "4"},
		{"Five", Five, "5"},
		{"Six", Six, "6"},
		{"Seven", Seven, "7"},
		{"Eight", Eight, "8"},
		{"Nine", Nine, "9"},
		{"Ten", Ten, "10"},
		{"Jack", Jack, "Jack"},
		{"Queen", Queen, "Queen"},
		{"King", King, "King"},
		{"Ace", Ace, "Ace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.rank) != tt.expected {
				t.Errorf("Rank %s = %s, want %s", tt.name, string(tt.rank), tt.expected)
			}
		})
	}
}

func TestCards_AddCard(t *testing.T) {
	t.Run("Add card to empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{},
		}

		newCard := Card{Rank: Ace, Suit: Spades}
		cards.Cards = append(cards.Cards, newCard)

		if len(cards.Cards) != 1 {
			t.Errorf("Expected deck length 1, got %d", len(cards.Cards))
		}

		if !reflect.DeepEqual(cards.Cards[0], newCard) {
			t.Errorf("Added card = %v, want %v", cards.Cards[0], newCard)
		}
	})

	t.Run("Add multiple cards", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{},
		}

		cardsToAdd := []Card{
			{Rank: Ace, Suit: Spades},
			{Rank: King, Suit: Hearts},
			{Rank: Queen, Suit: Clubs},
		}

		for _, card := range cardsToAdd {
			cards.Cards = append(cards.Cards, card)
		}

		if len(cards.Cards) != len(cardsToAdd) {
			t.Errorf("Expected deck length %d, got %d", len(cardsToAdd), len(cards.Cards))
		}

		for i, expectedCard := range cardsToAdd {
			if !reflect.DeepEqual(cards.Cards[i], expectedCard) {
				t.Errorf("Card at index %d = %v, want %v", i, cards.Cards[i], expectedCard)
			}
		}
	})
}

func TestCards_IsEmpty(t *testing.T) {
	t.Run("Empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{},
		}

		if len(cards.Cards) != 0 {
			t.Errorf("Expected empty deck, got length %d", len(cards.Cards))
		}
	})

	t.Run("Non-empty deck", func(t *testing.T) {
		cards := &Cards{
			Cards: []Card{
				{Rank: Ace, Suit: Spades},
			},
		}

		if len(cards.Cards) == 0 {
			t.Errorf("Expected non-empty deck, got empty deck")
		}
	})
}

func TestCards_Count(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected int
	}{
		{
			name:     "Empty deck",
			cards:    []Card{},
			expected: 0,
		},
		{
			name: "Single card",
			cards: []Card{
				{Rank: Ace, Suit: Spades},
			},
			expected: 1,
		},
		{
			name: "Multiple cards",
			cards: []Card{
				{Rank: Ace, Suit: Spades},
				{Rank: King, Suit: Hearts},
				{Rank: Queen, Suit: Clubs},
				{Rank: Jack, Suit: Diamonds},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cards := &Cards{
				Cards: tt.cards,
			}

			count := len(cards.Cards)
			if count != tt.expected {
				t.Errorf("Card count = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestCard_String(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		expected string
	}{
		{
			name:     "Ace of Spades",
			card:     Card{Rank: Ace, Suit: Spades},
			expected: "Ace of Spades",
		},
		{
			name:     "King of Hearts",
			card:     Card{Rank: King, Suit: Hearts},
			expected: "King of Hearts",
		},
		{
			name:     "Two of Clubs",
			card:     Card{Rank: Two, Suit: Clubs},
			expected: "2 of Clubs",
		},
		{
			name:     "Ten of Diamonds",
			card:     Card{Rank: Ten, Suit: Diamonds},
			expected: "10 of Diamonds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tt.card.Rank) + " of " + string(tt.card.Suit)
			if result != tt.expected {
				t.Errorf("Card string = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestCards_Peek(t *testing.T) {
	t.Run("Peek at top card without dealing", func(t *testing.T) {
		expectedCard := Card{Rank: Ace, Suit: Spades}
		cards := &Cards{
			Cards: []Card{
				expectedCard,
				{Rank: King, Suit: Hearts},
			},
		}

		originalLength := len(cards.Cards)

		// Peek at top card (simulated by accessing first element)
		if len(cards.Cards) > 0 {
			topCard := cards.Cards[0]
			if !reflect.DeepEqual(topCard, expectedCard) {
				t.Errorf("Top card = %v, want %v", topCard, expectedCard)
			}
		}

		// Verify deck length unchanged
		if len(cards.Cards) != originalLength {
			t.Errorf("Deck length changed after peek: got %d, want %d", len(cards.Cards), originalLength)
		}
	})
}

func TestCards_DealMultiple(t *testing.T) {
	t.Run("Deal multiple cards", func(t *testing.T) {
		originalCards := []Card{
			{Rank: Ace, Suit: Spades},
			{Rank: King, Suit: Hearts},
			{Rank: Queen, Suit: Clubs},
			{Rank: Jack, Suit: Diamonds},
		}

		cards := &Cards{
			Cards: make([]Card, len(originalCards)),
		}
		copy(cards.Cards, originalCards)

		dealtCards := []Card{}

		// Deal 3 cards
		for i := 0; i < 3; i++ {
			card, ok := cards.Deal()
			if !ok {
				t.Errorf("Deal %d returned false", i+1)
			}
			dealtCards = append(dealtCards, card)
		}

		// Verify dealt cards match original order
		for i := 0; i < 3; i++ {
			if !reflect.DeepEqual(dealtCards[i], originalCards[i]) {
				t.Errorf("Dealt card %d = %v, want %v", i, dealtCards[i], originalCards[i])
			}
		}

		// Verify remaining card
		if len(cards.Cards) != 1 {
			t.Errorf("Remaining cards = %d, want 1", len(cards.Cards))
		}

		if !reflect.DeepEqual(cards.Cards[0], originalCards[3]) {
			t.Errorf("Remaining card = %v, want %v", cards.Cards[0], originalCards[3])
		}
	})
}
