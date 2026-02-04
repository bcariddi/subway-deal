package game

// GamePhase represents the overall game phase
type GamePhase string

const (
	PhaseSetup    GamePhase = "setup"
	PhasePlaying  GamePhase = "playing"
	PhaseFinished GamePhase = "finished"
)

// TurnPhase represents the current turn phase
type TurnPhase string

const (
	TurnPhaseDraw     TurnPhase = "draw"
	TurnPhaseActions  TurnPhase = "actions"
	TurnPhaseResponse TurnPhase = "response" // Waiting for target to respond (Fare Evasion)
	TurnPhaseDiscard  TurnPhase = "discard"
)

// PendingAction represents an action waiting for target player response
type PendingAction struct {
	Action          *Action  `json:"action"`
	SourcePlayerID  string   `json:"sourcePlayerId"`
	TargetPlayerIDs []string `json:"targetPlayerIds"`
	RespondedIDs    []string `json:"respondedIds"`    // Players who accepted/paid
	RentMultiplier  int      `json:"rentMultiplier"`  // For Rush Hour (1 = normal, 2 = doubled)
	RentColor       string   `json:"rentColor"`       // Color being charged rent on
	RentAmount      int      `json:"rentAmount"`      // Calculated rent amount
}

// GameState holds the complete state of a game
type GameState struct {
	ID                    string         `json:"id"`
	Players               []*Player      `json:"players"`
	Deck                  []Card         `json:"-"` // Hidden from clients
	DiscardPile           []Card         `json:"discardPile"`
	CurrentPlayerIndex    int            `json:"currentPlayerIndex"`
	Phase                 GamePhase      `json:"phase"`
	TurnPhase             TurnPhase      `json:"turnPhase"`
	ActionsPlayedThisTurn int            `json:"actionsPlayedThisTurn"`
	MaxActionsPerTurn     int            `json:"maxActionsPerTurn"`
	Winner                *Player        `json:"winner,omitempty"`
	PendingAction         *PendingAction `json:"pendingAction,omitempty"`
}

func NewGameState(id string) *GameState {
	return &GameState{
		ID:                id,
		Players:           make([]*Player, 0),
		Deck:              make([]Card, 0),
		DiscardPile:       make([]Card, 0),
		Phase:             PhaseSetup,
		TurnPhase:         TurnPhaseDraw,
		MaxActionsPerTurn: 3,
	}
}

func (gs *GameState) GetCurrentPlayer() *Player {
	if len(gs.Players) == 0 {
		return nil
	}
	return gs.Players[gs.CurrentPlayerIndex]
}

func (gs *GameState) GetPlayer(playerID string) *Player {
	for _, p := range gs.Players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
}

func (gs *GameState) NextPlayer() {
	gs.CurrentPlayerIndex = (gs.CurrentPlayerIndex + 1) % len(gs.Players)
	gs.ActionsPlayedThisTurn = 0
	gs.TurnPhase = TurnPhaseDraw
}

func (gs *GameState) IsGameOver() bool {
	return gs.Phase == PhaseFinished
}

func (gs *GameState) CheckWinCondition() bool {
	for _, player := range gs.Players {
		if player.HasWon() {
			gs.Winner = player
			gs.Phase = PhaseFinished
			return true
		}
	}
	return false
}

// HasPendingAction returns true if there's an action awaiting response
func (gs *GameState) HasPendingAction() bool {
	return gs.PendingAction != nil
}

// GetPendingTargets returns player IDs who haven't responded yet
func (gs *GameState) GetPendingTargets() []string {
	if gs.PendingAction == nil {
		return nil
	}
	pending := make([]string, 0)
	for _, targetID := range gs.PendingAction.TargetPlayerIDs {
		responded := false
		for _, respID := range gs.PendingAction.RespondedIDs {
			if respID == targetID {
				responded = true
				break
			}
		}
		if !responded {
			pending = append(pending, targetID)
		}
	}
	return pending
}

// MarkResponded marks a player as having responded to the pending action
func (gs *GameState) MarkResponded(playerID string) {
	if gs.PendingAction == nil {
		return
	}
	gs.PendingAction.RespondedIDs = append(gs.PendingAction.RespondedIDs, playerID)
}

// AllTargetsResponded returns true if all targets have responded
func (gs *GameState) AllTargetsResponded() bool {
	return len(gs.GetPendingTargets()) == 0
}

// ClearPendingAction removes the pending action and returns to action phase
func (gs *GameState) ClearPendingAction() {
	gs.PendingAction = nil
	gs.TurnPhase = TurnPhaseActions
}
