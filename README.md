# Subway Deal

A fast-paced property trading card game themed around the NYC subway system. Collect subway line properties, charge rent, and be the first to complete 3 full sets to win.

## How to Play

### Objective
Be the first player to collect **3 complete property sets** (matching subway lines).

### Turn Structure
1. **Draw 2 cards** from the deck at the start of your turn
2. **Play up to 3 cards** per turn:
   - Play property cards to build your collection
   - Bank money cards for paying rent
   - Use action cards against opponents
   - Charge rent on your property sets
3. **End your turn** - discard down to 7 cards if over the limit

### Card Types
- **Property Cards** - Subway stations grouped by line color (complete a set to score)
- **Wildcard Properties** - Can count as either of two colors (your choice)
- **Money Cards** - MetroCards of various denominations for your bank
- **Rent Cards** - Charge other players rent based on your properties
- **Action Cards** - Special abilities like stealing properties or drawing extra cards

### Key Rules
- **No change given** - If you overpay rent, you don't get change back
- **Fare Evasion** - A special card that blocks any action played against you
- **Complete sets are protected** - Most steal actions only work on incomplete sets

## Running the Game

### Prerequisites
- Go 1.21+ installed

### Web Version (Hot-Seat Multiplayer)
```bash
cd server
go build -o webserver ./cmd/web
./webserver
```
Then open http://localhost:8080 in your browser.

### CLI Version
```bash
cd server
go build -o cli ./cmd/cli
./cli
```

## Project Structure

```
subway-deal/
├── server/
│   ├── game/           # Core game engine (Go)
│   │   ├── cards.go    # Card type definitions
│   │   ├── deck.go     # Deck creation and shuffling
│   │   ├── engine.go   # Game logic and turn management
│   │   ├── executor.go # Action execution and resolution
│   │   ├── player.go   # Player state and methods
│   │   └── state.go    # Game state management
│   ├── cmd/
│   │   ├── cli/        # Command-line interface
│   │   └── web/        # HTTP server for web UI
│   └── web/            # Frontend (HTML/CSS/JS)
│       ├── index.html
│       ├── styles.css
│       └── game.js
├── data/               # Card data (JSON)
└── CLAUDE.md           # Development notes
```

## MTA Line Colors

The game uses official MTA signage colors:

| Color | Lines | Hex |
|-------|-------|-----|
| Brown | J/Z | #996633 |
| Blue | A/C/E | #0039A6 |
| Orange | B/D/F/M | #FF6319 |
| Red | 1/2/3 | #EE352E |
| Yellow | N/Q/R/W | #FCCC0A |
| Green | 4/5/6 | #00933C |
| Gray | Shuttles | #808183 |
| Dark Blue | 7 | #2A344D |
| Light Green | G | #6CBE45 |
| Black | Regional Rail | #000000 |

## License

This is a fan project for personal/educational use.
