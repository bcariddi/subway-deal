package game

import "time"

// ActionType represents the type of action a player can take
type ActionType string

const (
	// Turn structure
	ActionDrawCards ActionType = "DRAW_CARDS"
	ActionEndTurn   ActionType = "END_TURN"

	// Playing cards
	ActionPlayProperty ActionType = "PLAY_PROPERTY"
	ActionPlayMoney    ActionType = "PLAY_MONEY"
	ActionPlayRent     ActionType = "PLAY_RENT"
	ActionPlayAction   ActionType = "PLAY_ACTION"
	ActionFlipWildcard ActionType = "FLIP_WILDCARD"

	// Specific action card effects
	ActionSwipeIn       ActionType = "SWIPE_IN"
	ActionFareEvasion   ActionType = "FARE_EVASION"
	ActionPowerBroker   ActionType = "POWER_BROKER"
	ActionServiceChange ActionType = "SERVICE_CHANGE"
	ActionLineClosure   ActionType = "LINE_CLOSURE"
	ActionMissedTrain   ActionType = "MISSED_YOUR_TRAIN"
	ActionItsMyStop     ActionType = "ITS_MY_STOP"
	ActionRushHour      ActionType = "RUSH_HOUR"
	ActionExpressService ActionType = "EXPRESS_SERVICE"
	ActionNewStation    ActionType = "NEW_STATION"

	// Responses to pending actions
	ActionAccept      ActionType = "ACCEPT"       // Accept/pay a pending action
	ActionPlayFareEvasion ActionType = "PLAY_FARE_EVASION" // Cancel with Fare Evasion
)

// Action represents a player action
type Action struct {
	Type      ActionType        `json:"type"`
	PlayerID  string            `json:"playerId"`
	Data      map[string]string `json:"data"` // Flexible key-value data
	Timestamp time.Time         `json:"timestamp"`
}

func NewAction(actionType ActionType, playerID string, data map[string]string) *Action {
	return &Action{
		Type:      actionType,
		PlayerID:  playerID,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Helper methods to get typed data from action
func (a *Action) GetCardID() string {
	return a.Data["cardId"]
}

func (a *Action) GetTargetPlayerID() string {
	return a.Data["targetPlayerId"]
}

func (a *Action) GetColor() string {
	return a.Data["color"]
}

func (a *Action) GetTargetCardID() string {
	return a.Data["targetCardId"]
}
