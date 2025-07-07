package deck

type Deck struct {
	Cards
}

func New() Deck {
	deck := Deck{
		Cards: Cards{Cards: []Card{}},
	}

	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck.Cards.Cards = append(deck.Cards.Cards, Card{Rank: rank, Suit: suit})
		}
	}

	return deck
}
