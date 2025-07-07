package blackjack

import "github.com/rbrabson/blackjack/deck"

type Shoe struct {
	reshufflePercentage int
	numDecks            int
	deck.Cards
}

func NewShoe(numDecks int, reshufflePercentage int) Shoe {
	shoe := Shoe{
		reshufflePercentage: reshufflePercentage,
		numDecks:            numDecks,
	}
	shoe.ReShuffle()

	return shoe
}

func (s *Shoe) StartNewRound() {
	cut := (s.reshufflePercentage * len(s.Cards.Cards)) / 100
	if len(s.Cards.Cards) <= cut {
		s.ReShuffle()
	}
}

func (s *Shoe) ReShuffle() {
	s.Cards = deck.Cards{Cards: []deck.Card{}}
	for i := 0; i < s.numDecks; i++ {
		deck := deck.New()
		s.Cards.Cards = append(s.Cards.Cards, deck.Cards.Cards...)
	}
	s.Shuffle()
}
