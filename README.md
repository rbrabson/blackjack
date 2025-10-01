# Blackjack Game

A comprehensive console-based blackjack game written in Go, featuring a dealer, multiple players, and a multi-deck shoe.

## Features

- **Multiple Players**: Support for one or more players with individual chip counts and betting
- **Professional Dealer**: AI dealer that follows standard blackjack rules (hits on soft 17)
- **Multi-Deck Shoe**: Configurable shoe with multiple decks, automatic shuffling, and cut card placement
- **Standard Blackjack Rules**:
  - Blackjack pays 3:2
  - Dealer hits on soft 17
  - Double down on any two cards
  - Split pairs (including multiple splits)
  - Standard hand evaluation with proper Ace handling
- **Interactive Gameplay**: Console-based interface with clear prompts and game state display
- **Betting System**: Chip-based betting with win/loss tracking

## Game Components

### 🃏 Hand

- Tracks cards and calculates blackjack values
- Handles Ace as 1 or 11 automatically
- Detects blackjack, busts, and soft hands
- Provides both visible and hidden display modes

### 👤 Player

- Manages multiple hands (for splits), chips, and bets
- Supports hit, stand, double down, and split actions
- Tracks active/inactive status during rounds
- Handles win/loss payouts

### 🎯 Dealer

- Follows standard blackjack dealer rules
- Automatically hits on 16 or less, stands on 17 or more
- Hits on soft 17 (configurable house rule)
- Manages hole card display

### 🎴 Shoe

- Multi-deck shoe with configurable deck count (default: 6 decks)
- Automatic shuffling with cut card placement
- Tracks penetration percentage
- Reshuffles when cut card is reached

### 🎮 Game Engine

- Orchestrates complete blackjack rounds
- Handles betting, dealing, player actions, and payouts
- Tracks game statistics and round progression
- Manages game state and player turns

## How to Play

1. **Setup**: Add players with starting chip amounts
2. **Betting**: Each player places their bet for the round
3. **Dealing**: Two cards dealt to each player and dealer (dealer's second card is face down)
4. **Player Actions**: Each player can:
   - **Hit**: Take another card
   - **Stand**: Keep current hand
   - **Double Down**: Double the bet and take exactly one more card (if eligible)
   - **Split**: If dealt a pair, split into two separate hands with separate bets
5. **Dealer Play**: Dealer reveals hole card and hits according to rules
6. **Payouts**: Winners are paid according to standard blackjack payouts

## Running the Game

```bash
# Build the game
go build

# Run the game
./blackjack
```

## Game Rules

- **Blackjack**: 21 with first two cards (pays 3:2)
- **Bust**: Hand value over 21 (automatic loss)
- **Dealer Rules**: Must hit on 16 or less, stand on 17 or more
- **Soft 17**: Dealer hits on soft 17 (Ace + 6)
- **Double Down**: Available on any two cards if you have sufficient chips
- **Split**: Available when dealt a pair (two cards of same rank)
  - Each split hand gets a separate bet equal to the original bet
  - Split hands cannot achieve "natural" blackjack (still pays 1:1)
  - Can continue to hit, stand, or double down on each split hand
- **Winning**: Beat dealer without busting, or dealer busts

## Dependencies

- `github.com/rbrabson/cards`: Provides card, deck, and shoe functionality

## File Structure

- `cmd/blackjack/main.go`: Main game loop and user interface
- `game.go`: Core game logic and round management
- `hand.go`: Hand representation and value calculation
- `player.go`: Player management and actions
- `dealer.go`: Dealer logic and rules
- `shoe.go`: Multi-deck shoe with blackjack-specific features

## Example Gameplay

```sh
🃏 Welcome to Blackjack! 🃏
========================

Enter player name (or 'done' to start): Alice
Enter starting chips: 1000
Added Alice with 1000 chips.

Enter player name (or 'done' to start): Bob
Enter starting chips: 500
Added Bob with 500 chips.

Enter player name (or 'done' to start): done

🎲 Starting Round 1 🎲
===================

Alice (Chips: 1000), place your bet: 50
Alice bet 50 chips.

Bob (Chips: 500), place your bet: 25
Bob bet 25 chips.

📋 Initial Cards:
Dealer: [Hidden, 7♠] (Visible Value: 7)

Alice (Chips: 950, Bet: 50, active): [K♥, 5♦] (Value: 15)
Bob (Chips: 475, Bet: 25, active): [A♠, Q♣] (Value: 21)

🎯 Bob has blackjack!

🎮 Alice's turn:
Alice: [K♥, 5♦] (Value: 15)
Choose action: (h)it, (s)tand, (d)ouble down: h
Drew: [K♥, 5♦, 8♣] (Value: 23)
💥 Alice busted!

🎯 Dealer's turn:
Revealing hole card...
Dealer: [6♥, 7♠] (Value: 13)
Dealer: [6♥, 7♠, 9♦] (Value: 22)

🏁 Final Results:
Dealer: [6♥, 7♠, 9♦] (Value: 22)

Alice (Chips: 950, Bet: 50, inactive): [K♥, 5♦, 8♣] (Value: 23)
Bob (Chips: 475, Bet: 25, inactive): [A♠, Q♣] (Value: 21)

💰 Round Results:
================
Alice: Dealer Wins
  Chips: 950
Bob: Player Blackjack!
  Chips: 512

Play another round? (y/n):
```

## Split Example

```
🎮 Charlie's turn:
Charlie: [8♠, 8♥] (Value: 16)
Choose action: (h)it, (s)tand, (d)ouble down, s(p)lit: p
Hand split! You now have 2 hands.
Current hand: [8♠, 3♦] (Value: 11)

Charlie - Hand 1 of 2: [8♠, 3♦] (Value: 11) (Split)
Choose action: (h)it, (s)tand, (d)ouble down: h
Drew: [8♠, 3♦, 7♣] (Value: 18) (Split)
Standing on hand.

Charlie - Hand 2 of 2: [8♥, K♠] (Value: 18) (Split)
Choose action: (h)it, (s)tand: s
Standing on hand.

✅ Charlie finished all hands.

💰 Round Results:
================
Charlie:
  Hand 1: Player Wins
  Hand 2: Player Wins
  Final Chips: 1100
```

Enjoy playing blackjack! 🎉
