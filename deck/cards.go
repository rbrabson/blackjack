package deck

import "math/rand/v2"

// Suit represents the suit of a card in a standard deck.
type Suit string

const (
	Clubs    Suit = "Clubs"
	Diamonds Suit = "Diamonds"
	Hearts   Suit = "Hearts"
	Spades   Suit = "Spades"
)

// Rank represents the rank of a card in a standard deck.
type Rank string

const (
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "Jack"
	Queen Rank = "Queen"
	King  Rank = "King"
	Ace   Rank = "Ace"
)

// Card represents a single playing card with a rank and suit.
type Card struct {
	Rank Rank
	Suit Suit
}

// Cards represents a collection of playing cards.
type Cards struct {
	Cards []Card
}

// Shuffle randomly shuffles the cards in the collection.
func (c *Cards) Shuffle() {
	rand.Shuffle(len(c.Cards), func(i, j int) {
		c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i]
	})
}

// Deal removes and returns the top card from the collection.
func (c *Cards) Deal() (Card, bool) {
	if len(c.Cards) == 0 {
		return Card{}, false
	}
	card := c.Cards[0]
	c.Cards = c.Cards[1:]
	return card, true
}
