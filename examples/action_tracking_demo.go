package main

import (
	"fmt"

	"github.com/rbrabson/blackjack"
	"github.com/rbrabson/cards"
)

func demonstrateActionTracking() {
	fmt.Println("🎯 Action Tracking Demonstration")
	fmt.Println("================================")

	// Create a player and dealer
	player := blackjack.NewPlayer("Alice", 1000)
	dealer := blackjack.NewDealer()

	// Simulate initial deal
	fmt.Println("\n📋 Initial Deal:")
	player.DealCard(cards.Card{Suit: cards.Spades, Rank: cards.Ten})
	player.DealCard(cards.Card{Suit: cards.Hearts, Rank: cards.Six})
	dealer.DealCard(cards.Card{Suit: cards.Clubs, Rank: cards.King})
	dealer.DealCard(cards.Card{Suit: cards.Diamonds, Rank: cards.Seven})

	fmt.Printf("Player Hand: %s\n", player.CurrentHand().String())
	fmt.Printf("Player Actions: %s\n", player.CurrentHand().ActionSummary())
	fmt.Printf("Dealer Hand: %s\n", dealer.Hand().String())
	fmt.Printf("Dealer Actions: %s\n", dealer.Hand().ActionSummary())

	// Player hits
	fmt.Println("\n🃏 Player Hits:")
	player.Hit(cards.Card{Suit: cards.Spades, Rank: cards.Five})
	fmt.Printf("Player Hand: %s\n", player.CurrentHand().String())
	fmt.Printf("Player Actions: %s\n", player.CurrentHand().ActionSummary())

	// Player stands
	fmt.Println("\n✋ Player Stands:")
	player.CurrentHand().Stand()
	fmt.Printf("Player Actions: %s\n", player.CurrentHand().ActionSummary())

	// Dealer plays
	fmt.Println("\n🎲 Dealer Plays:")
	dealer.Hit(cards.Card{Suit: cards.Hearts, Rank: cards.Two})
	dealer.Hit(cards.Card{Suit: cards.Clubs, Rank: cards.Three})
	dealer.Stand()
	fmt.Printf("Dealer Hand: %s\n", dealer.Hand().String())
	fmt.Printf("Dealer Actions: %s\n", dealer.Hand().ActionSummary())

	// Show detailed action history
	fmt.Println("\n📝 Detailed Action History:")
	fmt.Println("Player:")
	for i, action := range player.CurrentHand().Actions() {
		fmt.Printf("  %d. %s", i+1, action.Type)
		if action.Card != nil {
			fmt.Printf(" (%s)", action.Card.String())
		}
		if action.Details != "" {
			fmt.Printf(" - %s", action.Details)
		}
		fmt.Printf(" at %s\n", action.Timestamp.Format("15:04:05.000"))
	}

	fmt.Println("\nDealer:")
	for i, action := range dealer.Hand().Actions() {
		fmt.Printf("  %d. %s", i+1, action.Type)
		if action.Card != nil {
			fmt.Printf(" (%s)", action.Card.String())
		}
		if action.Details != "" {
			fmt.Printf(" - %s", action.Details)
		}
		fmt.Printf(" at %s\n", action.Timestamp.Format("15:04:05.000"))
	}

	// Demonstrate other actions
	fmt.Println("\n🎯 Other Action Examples:")

	// Double Down example
	fmt.Println("\n💰 Double Down Example:")
	player2 := blackjack.NewPlayer("Bob", 1000)
	player2.DealCard(cards.Card{Suit: cards.Spades, Rank: cards.Nine})
	player2.DealCard(cards.Card{Suit: cards.Hearts, Rank: cards.Two})
	player2.PlaceBet(50)

	player2.DoubleDown()
	player2.DoubleDownHit(cards.Card{Suit: cards.Clubs, Rank: cards.Ten})
	fmt.Printf("Bob's Hand: %s\n", player2.CurrentHand().String())
	fmt.Printf("Bob's Actions: %s\n", player2.CurrentHand().ActionSummary())

	// Surrender example
	fmt.Println("\n🏳️ Surrender Example:")
	player3 := blackjack.NewPlayer("Charlie", 1000)
	player3.DealCard(cards.Card{Suit: cards.Spades, Rank: cards.Ten})
	player3.DealCard(cards.Card{Suit: cards.Hearts, Rank: cards.Six})
	player3.PlaceBet(100)

	player3.Surrender()
	fmt.Printf("Charlie's Hand: %s\n", player3.CurrentHand().String())
	fmt.Printf("Charlie's Actions: %s\n", player3.CurrentHand().ActionSummary())

	// Split example
	fmt.Println("\n✂️ Split Example:")
	player4 := blackjack.NewPlayer("Diana", 1000)
	player4.DealCard(cards.Card{Suit: cards.Spades, Rank: cards.Eight})
	player4.DealCard(cards.Card{Suit: cards.Hearts, Rank: cards.Eight})
	player4.PlaceBet(75)

	player4.Split()
	fmt.Printf("Diana's Hand 1: %s\n", player4.Hands()[0].String())
	fmt.Printf("Diana's Hand 1 Actions: %s\n", player4.Hands()[0].ActionSummary())
	fmt.Printf("Diana's Hand 2: %s\n", player4.Hands()[1].String())
	fmt.Printf("Diana's Hand 2 Actions: %s\n", player4.Hands()[1].ActionSummary())

	fmt.Println("\n✅ Action tracking demonstration complete!")
}
