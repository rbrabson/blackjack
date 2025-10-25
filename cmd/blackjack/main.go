package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rbrabson/blackjack"
)

func main() {
	fmt.Println("🃏 Welcome to Blackjack! 🃏")
	fmt.Println("========================")

	// Create a new game with 6 decks (typical casino setup)
	game := blackjack.New(6)

	// Setup players
	setupPlayers(game)

	// Main game loop
	for {
		if !playRound(game) {
			break
		}

		// Check if any players want to continue
		if !askToContinue(game) {
			break
		}
	}

	fmt.Println("\n🎉 Thanks for playing Blackjack! 🎉")
	showFinalStats(game)
}

func setupPlayers(game *blackjack.Game) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nEnter player name (or 'done' to start): ")
		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())

		if strings.ToLower(name) == "done" {
			break
		}

		if name == "" {
			fmt.Println("Please enter a valid name.")
			continue
		}

		// Check if player already exists
		if game.GetPlayer(name) != nil {
			fmt.Println("Player with that name already exists.")
			continue
		}

		fmt.Print("Enter starting chips: ")
		scanner.Scan()
		chipsStr := strings.TrimSpace(scanner.Text())
		chips, err := strconv.Atoi(chipsStr)
		if err != nil || chips <= 0 {
			fmt.Println("Please enter a valid positive number for chips.")
			continue
		}

		game.AddPlayer(name, blackjack.WithChips(chips))
		fmt.Printf("Added %s with %d chips.\n", name, chips)
	}

	if len(game.Players()) == 0 {
		// Add a default player if none were added
		fmt.Println("No players added. Adding default player 'Player1' with 1000 chips.")
		game.AddPlayer("Player1", blackjack.WithChips(1000))
	}
}

func playRound(game *blackjack.Game) bool {
	// Start new round
	err := game.StartNewRound()
	if err != nil {
		fmt.Printf("Error starting round: %v\n", err)
		return false
	}

	fmt.Printf("\n🎲 Starting Round %d 🎲\n", game.Round())
	fmt.Println("===================")

	// Place bets
	if !placeBets(game) {
		return false
	}

	// Deal initial cards
	err = game.DealInitialCards()
	if err != nil {
		fmt.Printf("Error dealing cards: %v\n", err)
		return false
	}

	// Show initial game state
	fmt.Println("\n📋 Initial Cards:")
	fmt.Println(game.GetGameStatus(false))

	// Check for dealer blackjack
	if game.Dealer().HasBlackjack() {
		fmt.Println("🎯 Dealer has blackjack!")
		fmt.Println(game.GetGameStatus(true))
		game.PayoutResults()
		showRoundResults(game)
		return true
	}

	// Player turns
	playerTurns(game)

	// Dealer turn (if any players are still in)
	if hasActiveNonBustedPlayers(game) {
		fmt.Println("\n🎯 Dealer's turn:")
		fmt.Println("Revealing hole card...")
		fmt.Println(game.Dealer().RevealHoleCard())

		err = game.DealerPlay()
		if err != nil {
			fmt.Printf("Error during dealer play: %v\n", err)
			return false
		}

		fmt.Println("\nDealer finished:")
		fmt.Println(game.Dealer().String())
	}

	// Show final results
	fmt.Println("\n🏁 Final Results:")
	fmt.Println(game.GetGameStatus(true))

	// Pay out results
	game.PayoutResults()
	showRoundResults(game)

	return true
}

func placeBets(game *blackjack.Game) bool {
	scanner := bufio.NewScanner(os.Stdin)

	for _, player := range game.Players() {
		if player.Chips() <= 0 {
			fmt.Printf("%s has no chips left and will sit out this round.\n", player.Name())
			player.SetActive(false)
			continue
		}

		for {
			fmt.Printf("\n%s (Chips: %d), place your bet: ", player.Name(), player.Chips())
			scanner.Scan()
			betStr := strings.TrimSpace(scanner.Text())

			if betStr == "quit" {
				return false
			}

			bet, err := strconv.Atoi(betStr)
			if err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}

			hand := player.CurrentHand()
			err = hand.PlaceBet(bet)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Printf("%s bet %d chips.\n", player.Name(), bet)
			break
		}
	}

	// Check if any players placed bets
	hasActivePlayers := false
	for _, player := range game.Players() {
		hand := player.CurrentHand()
		if player.IsActive() && hand.Bet() > 0 {
			hasActivePlayers = true
			break
		}
	}

	return hasActivePlayers
}

func playerTurns(game *blackjack.Game) {
	scanner := bufio.NewScanner(os.Stdin)

	for _, player := range game.Players() {
		hand := player.CurrentHand()
		if !player.IsActive() || hand.Bet() == 0 {
			continue
		}

		fmt.Printf("\n🎮 %s's turn:\n", player.Name())

		// Handle all hands for this player (including splits)
		for player.HasActiveHands() {
			currentHand := player.CurrentHand()

			// Check for player blackjack
			if currentHand.IsBlackjack() {
				fmt.Printf("🎯 %s has blackjack on hand %d!\n", player.Name(), player.GetCurrentHandIndex()+1)
				if !player.MoveToNextActiveHand() {
					player.SetActive(false)
					break
				}
				continue
			}

			// Show current hand status
			if len(player.Hands()) > 1 {
				fmt.Printf("\n%s - Hand %d of %d: %s\n",
					player.Name(),
					player.GetCurrentHandIndex()+1,
					len(player.Hands()),
					currentHand.String())
			} else {
				fmt.Printf("\n%s: %s\n", player.Name(), currentHand.String())
			}

			// Player actions for current hand
			for currentHand.IsActive() && !currentHand.IsBusted() && !currentHand.IsBlackjack() {
				fmt.Print("Choose action: (h)it, (s)tand")

				if currentHand.CanDoubleDown() {
					fmt.Print(", (d)ouble down")
				}

				if currentHand.CanSplit() {
					fmt.Print(", s(p)lit")
				}

				if currentHand.CanSurrender() {
					fmt.Print(", s(u)rrender")
				}

				fmt.Print(": ")
				scanner.Scan()
				action := strings.ToLower(strings.TrimSpace(scanner.Text()))

				switch action {
				case "h", "hit":
					err := game.PlayerHit(player.Name())
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					fmt.Printf("Drew: %s\n", currentHand.String())

					if currentHand.IsBusted() {
						fmt.Printf("💥 Hand busted!\n")
						currentHand.SetActive(false)
					}

				case "s", "stand":
					fmt.Printf("Standing on hand.\n")
					err := game.PlayerStand(player.Name())
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

				case "d", "double", "double down":
					if !currentHand.CanDoubleDown() {
						fmt.Println("Cannot double down.")
						continue
					}

					err := currentHand.DoubleDown()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					err = game.PlayerDoubleDownHit(player.Name())
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					fmt.Printf("Doubled down! Drew: %s\n", currentHand.String())

					if currentHand.IsBusted() {
						fmt.Printf("💥 Hand busted!\n")
					}

					// Double down ends the hand
					err = game.PlayerStand(player.Name())
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}

				case "p", "split":
					if !currentHand.CanSplit() {
						fmt.Println("Cannot split.")
						continue
					}

					err := currentHand.Split()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					fmt.Printf("Hand split! You now have %d hands.\n", len(player.Hands()))
					// Show current hand after split
					fmt.Printf("Current hand: %s\n", currentHand.String())

				case "u", "surrender":
					if !currentHand.CanSurrender() {
						fmt.Println("Cannot surrender.")
						continue
					}

					err := game.PlayerSurrender(player.Name())
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						continue
					}

					fmt.Printf("Surrendered! Half bet returned.\n")

				default:
					fmt.Println("Invalid action. Please choose (h)it, (s)tand, (d)ouble down, s(p)lit, or s(u)rrender if available.")
				}
			}

			// Move to next hand if current hand is done
			if !currentHand.IsActive() {
				if !player.MoveToNextActiveHand() {
					player.SetActive(false)
					break
				}
			}
		}

		fmt.Printf("✅ %s finished all hands.\n", player.Name())
	}
}

func hasActiveNonBustedPlayers(game *blackjack.Game) bool {
	for _, player := range game.Players() {
		hand := player.CurrentHand()
		if hand.Bet() > 0 && !hand.IsBusted() {
			return true
		}
	}
	return false
}

func showRoundResults(game *blackjack.Game) {
	fmt.Println("\n💰 Round Results:")
	fmt.Println("================")

	for _, player := range game.Players() {
		hands := player.Hands()
		if len(hands) == 1 {
			// Single hand
			result := game.EvaluateHand(player.CurrentHand())
			fmt.Printf("%s: %s\n", player.Name(), result.String())
		} else {
			// Multiple hands (splits)
			fmt.Printf("%s:\n", player.Name())
			for idx, hand := range hands {
				// Temporarily set current hand for evaluation
				result := game.EvaluateHand(hand)
				fmt.Printf("  Hand %d: %s\n", idx+1, result.String())
			}
		}

		fmt.Printf("  Final Chips: %d\n", player.Chips())
	}
}

func askToContinue(game *blackjack.Game) bool {
	// Check if any players have chips left
	hasChips := false
	for _, player := range game.Players() {
		if player.Chips() > 0 {
			hasChips = true
			break
		}
	}

	if !hasChips {
		fmt.Println("\nNo players have chips remaining. Game over!")
		return false
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nPlay another round? (y/n): ")
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	return response == "y" || response == "yes"
}

func showFinalStats(game *blackjack.Game) {
	fmt.Println("\n📊 Final Statistics:")
	fmt.Println("===================")
	fmt.Printf("Rounds played: %d\n", game.Round())
	fmt.Printf("Shoe penetration: %.1f%%\n", game.Shoe().Penetration())

	fmt.Println("\nFinal chip counts:")
	for _, player := range game.Players() {
		fmt.Printf("  %s: %d chips\n", player.Name(), player.Chips())
	}
}
