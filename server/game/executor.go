package game

import (
	"fmt"
	"math/rand"
	"strconv"
)

// ActionResult contains the result of executing an action
type ActionResult struct {
	Success       bool            `json:"success"`
	Message       string          `json:"message"`
	Error         string          `json:"error,omitempty"`
	Payments      []PaymentResult `json:"payments,omitempty"`
	PendingAction bool            `json:"pendingAction,omitempty"` // True if waiting for responses
}

type PaymentResult struct {
	FromPlayer string `json:"fromPlayer"`
	Amount     int    `json:"amount"`
	Shortfall  int    `json:"shortfall"`
}

// Executor executes validated actions
type Executor struct {
	state *GameState
}

func NewExecutor(state *GameState) *Executor {
	return &Executor{state: state}
}

func (e *Executor) Execute(action *Action) ActionResult {
	switch action.Type {
	// Non-targetable actions (execute immediately)
	case ActionPlayProperty:
		return e.executePlayProperty(action)
	case ActionPlayMoney:
		return e.executePlayMoney(action)
	case ActionSwipeIn:
		return e.executeSwipeIn(action)
	case ActionExpressService:
		return e.executeExpressService(action)
	case ActionNewStation:
		return e.executeNewStation(action)
	case ActionFlipWildcard:
		return e.executeFlipWildcard(action)
	case ActionEndTurn:
		return ActionResult{Success: true, Message: "turn ended"}

	// Targetable actions (create pending action for Fare Evasion response)
	case ActionPlayRent:
		return e.createPendingRent(action)
	case ActionPowerBroker:
		return e.createPendingPowerBroker(action)
	case ActionLineClosure:
		return e.createPendingLineClosure(action)
	case ActionServiceChange:
		return e.createPendingServiceChange(action)
	case ActionMissedTrain:
		return e.createPendingMissedTrain(action)
	case ActionItsMyStop:
		return e.createPendingItsMyStop(action)

	// Response actions
	case ActionAccept:
		return e.executeAccept(action)
	case ActionPlayFareEvasion:
		return e.executeFareEvasion(action)

	default:
		return ActionResult{Success: false, Error: fmt.Sprintf("no executor for: %s", action.Type)}
	}
}

// ========== Non-targetable actions (immediate execution) ==========

func (e *Executor) executePlayProperty(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	var color string
	switch c := card.(type) {
	case *PropertyCard:
		color = c.Color
	case *WildcardCard:
		color = c.CurrentColor
	}

	propertySet := player.GetPropertySet(color)
	propertySet.AddCard(card)
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s played %s", player.Name, card.GetName()),
	}
}

func (e *Executor) executePlayMoney(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	player.AddToBank(card)
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s banked %s ($%d)", player.Name, card.GetName(), card.GetValue()),
	}
}

func (e *Executor) executeSwipeIn(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	for i := 0; i < 2; i++ {
		drawn := e.drawCard()
		if drawn != nil {
			player.AddToHand(drawn)
		}
	}

	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s used Swipe In and drew 2 cards", player.Name),
	}
}

func (e *Executor) executeExpressService(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	color := action.GetColor()
	propertySet := player.GetPropertySet(color)

	if err := propertySet.AddImprovement("express"); err != nil {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s added Express Service to %s set (+$3 rent)", player.Name, color),
	}
}

func (e *Executor) executeNewStation(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	color := action.GetColor()
	propertySet := player.GetPropertySet(color)

	if err := propertySet.AddImprovement("station"); err != nil {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s added New Station to %s set (+$4 rent)", player.Name, color),
	}
}

func (e *Executor) executeFlipWildcard(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	cardID := action.GetCardID()

	for color, ps := range player.Properties {
		for i, card := range ps.Cards {
			if card.GetID() == cardID {
				wc, ok := card.(*WildcardCard)
				if !ok {
					return ActionResult{Success: false, Error: "card is not a wildcard"}
				}

				ps.Cards = append(ps.Cards[:i], ps.Cards[i+1:]...)
				wc.FlipColor()
				player.GetPropertySet(wc.CurrentColor).AddCard(wc)

				return ActionResult{
					Success: true,
					Message: fmt.Sprintf("%s flipped wildcard from %s to %s", player.Name, color, wc.CurrentColor),
				}
			}
		}
	}

	return ActionResult{Success: false, Error: "wildcard not found"}
}

// ========== Targetable actions (create pending action) ==========

func (e *Executor) createPendingRent(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	rentCard := card.(*RentCard)
	color := action.GetColor()
	propertySet := player.GetPropertySet(color)
	rentAmount := propertySet.GetRent()

	// Check for Rush Hour modifier
	multiplier := 1
	if rushHourID := action.Data["rushHourCardId"]; rushHourID != "" {
		rushCard, err := player.RemoveFromHand(rushHourID)
		if err == nil {
			e.state.DiscardPile = append(e.state.DiscardPile, rushCard)
			multiplier = 2
		}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	// Determine targets
	var targetIDs []string
	if rentCard.Target == "all" {
		for _, p := range e.state.Players {
			if p.ID != player.ID {
				targetIDs = append(targetIDs, p.ID)
			}
		}
	} else {
		targetIDs = []string{action.GetTargetPlayerID()}
	}

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: targetIDs,
		RespondedIDs:    make([]string, 0),
		RentMultiplier:  multiplier,
		RentColor:       color,
		RentAmount:      rentAmount * multiplier,
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	msg := fmt.Sprintf("%s demands $%d rent on %s", player.Name, rentAmount*multiplier, color)
	if multiplier == 2 {
		msg += " (Rush Hour!)"
	}

	return ActionResult{
		Success:       true,
		Message:       msg,
		PendingAction: true,
	}
}

func (e *Executor) createPendingPowerBroker(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	targetPlayer := e.state.GetPlayer(action.GetTargetPlayerID())
	targetSet := targetPlayer.GetPropertySet(action.GetColor())

	if targetSet.IsComplete() {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: "cannot steal from complete set"}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: []string{action.GetTargetPlayerID()},
		RespondedIDs:    make([]string, 0),
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success:       true,
		Message:       fmt.Sprintf("%s plays Power Broker against %s", player.Name, targetPlayer.Name),
		PendingAction: true,
	}
}

func (e *Executor) createPendingLineClosure(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	targetPlayer := e.state.GetPlayer(action.GetTargetPlayerID())
	targetSet := targetPlayer.GetPropertySet(action.GetColor())

	if !targetSet.IsComplete() {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: "can only steal complete sets"}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: []string{action.GetTargetPlayerID()},
		RespondedIDs:    make([]string, 0),
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success:       true,
		Message:       fmt.Sprintf("%s plays Line Closure against %s's %s set!", player.Name, targetPlayer.Name, action.GetColor()),
		PendingAction: true,
	}
}

func (e *Executor) createPendingServiceChange(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	playerColor := action.Data["playerColor"]
	playerSet := player.GetPropertySet(playerColor)
	targetPlayer := e.state.GetPlayer(action.GetTargetPlayerID())
	targetSet := targetPlayer.GetPropertySet(action.GetColor())

	if playerSet.IsComplete() {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: "cannot swap from your complete set"}
	}
	if targetSet.IsComplete() {
		player.AddToHand(card)
		return ActionResult{Success: false, Error: "cannot swap from opponent's complete set"}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: []string{action.GetTargetPlayerID()},
		RespondedIDs:    make([]string, 0),
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success:       true,
		Message:       fmt.Sprintf("%s plays Service Change against %s", player.Name, targetPlayer.Name),
		PendingAction: true,
	}
}

func (e *Executor) createPendingMissedTrain(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	targetPlayer := e.state.GetPlayer(action.GetTargetPlayerID())

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: []string{action.GetTargetPlayerID()},
		RespondedIDs:    make([]string, 0),
		RentAmount:      5, // $5 debt
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success:       true,
		Message:       fmt.Sprintf("%s demands $5 from %s (Missed Your Train)", player.Name, targetPlayer.Name),
		PendingAction: true,
	}
}

func (e *Executor) createPendingItsMyStop(action *Action) ActionResult {
	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)

	var targetIDs []string
	for _, p := range e.state.Players {
		if p.ID != player.ID {
			targetIDs = append(targetIDs, p.ID)
		}
	}

	e.state.PendingAction = &PendingAction{
		Action:          action,
		SourcePlayerID:  player.ID,
		TargetPlayerIDs: targetIDs,
		RespondedIDs:    make([]string, 0),
		RentAmount:      2, // $2 from each
	}
	e.state.TurnPhase = TurnPhaseResponse
	e.state.ActionsPlayedThisTurn++

	return ActionResult{
		Success:       true,
		Message:       fmt.Sprintf("%s demands $2 from everyone (It's My Stop!)", player.Name),
		PendingAction: true,
	}
}

// ========== Response actions ==========

func (e *Executor) executeAccept(action *Action) ActionResult {
	pending := e.state.PendingAction
	if pending == nil {
		return ActionResult{Success: false, Error: "no pending action"}
	}

	responderID := action.PlayerID
	e.state.MarkResponded(responderID)

	result := e.resolveForPlayer(responderID)

	// Check if all targets have responded
	if e.state.AllTargetsResponded() {
		e.state.ClearPendingAction()
	}

	return result
}

func (e *Executor) executeFareEvasion(action *Action) ActionResult {
	pending := e.state.PendingAction
	if pending == nil {
		return ActionResult{Success: false, Error: "no pending action"}
	}

	player := e.state.GetPlayer(action.PlayerID)
	card, err := player.RemoveFromHand(action.GetCardID())
	if err != nil {
		return ActionResult{Success: false, Error: err.Error()}
	}

	e.state.DiscardPile = append(e.state.DiscardPile, card)
	e.state.MarkResponded(action.PlayerID)

	// If this was a single-target action, cancel the whole thing
	// If multi-target, just remove this player from targets
	if len(pending.TargetPlayerIDs) == 1 {
		e.state.ClearPendingAction()
		return ActionResult{
			Success: true,
			Message: fmt.Sprintf("%s blocked the action with Fare Evasion!", player.Name),
		}
	}

	// Multi-target: check if all responded
	if e.state.AllTargetsResponded() {
		e.state.ClearPendingAction()
	}

	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s blocked with Fare Evasion!", player.Name),
	}
}

// resolveForPlayer executes the pending action's effect on a specific player who accepted
func (e *Executor) resolveForPlayer(playerID string) ActionResult {
	pending := e.state.PendingAction
	sourcePlayer := e.state.GetPlayer(pending.SourcePlayerID)
	targetPlayer := e.state.GetPlayer(playerID)

	switch pending.Action.Type {
	case ActionPlayRent, ActionMissedTrain, ActionItsMyStop:
		// Pay money
		amount := pending.RentAmount
		paid, total := targetPlayer.PayMoney(amount)
		for _, paidCard := range paid {
			sourcePlayer.AddToBank(paidCard)
		}
		return ActionResult{
			Success: true,
			Message: fmt.Sprintf("%s paid $%d", targetPlayer.Name, total),
			Payments: []PaymentResult{{
				FromPlayer: targetPlayer.Name,
				Amount:     total,
				Shortfall:  max(0, amount-total),
			}},
		}

	case ActionPowerBroker:
		// Steal single property
		targetColor := pending.Action.GetColor()
		targetSet := targetPlayer.GetPropertySet(targetColor)
		stolenCard, err := targetSet.RemoveCard(pending.Action.GetTargetCardID())
		if err != nil {
			return ActionResult{Success: false, Error: "target card not found"}
		}

		var color string
		switch c := stolenCard.(type) {
		case *PropertyCard:
			color = c.Color
		case *WildcardCard:
			color = c.CurrentColor
		}
		sourcePlayer.GetPropertySet(color).AddCard(stolenCard)

		return ActionResult{
			Success: true,
			Message: fmt.Sprintf("%s stole %s from %s", sourcePlayer.Name, stolenCard.GetName(), targetPlayer.Name),
		}

	case ActionLineClosure:
		// Steal complete set
		targetColor := pending.Action.GetColor()
		targetSet := targetPlayer.GetPropertySet(targetColor)
		playerSet := sourcePlayer.GetPropertySet(targetColor)

		for _, c := range targetSet.Cards {
			playerSet.AddCard(c)
		}
		for _, imp := range targetSet.Improvements {
			playerSet.Improvements = append(playerSet.Improvements, imp)
		}
		targetSet.Cards = nil
		targetSet.Improvements = nil

		return ActionResult{
			Success: true,
			Message: fmt.Sprintf("%s stole complete %s set from %s", sourcePlayer.Name, targetColor, targetPlayer.Name),
		}

	case ActionServiceChange:
		// Swap properties
		playerCardID := pending.Action.Data["playerCardId"]
		playerColor := pending.Action.Data["playerColor"]
		targetCardID := pending.Action.GetTargetCardID()
		targetColor := pending.Action.GetColor()

		playerSet := sourcePlayer.GetPropertySet(playerColor)
		targetSet := targetPlayer.GetPropertySet(targetColor)

		playerCard, _ := playerSet.RemoveCard(playerCardID)
		targetCard, _ := targetSet.RemoveCard(targetCardID)

		var newPlayerColor, newTargetColor string
		switch c := targetCard.(type) {
		case *PropertyCard:
			newPlayerColor = c.Color
		case *WildcardCard:
			newPlayerColor = c.CurrentColor
		}
		switch c := playerCard.(type) {
		case *PropertyCard:
			newTargetColor = c.Color
		case *WildcardCard:
			newTargetColor = c.CurrentColor
		}

		sourcePlayer.GetPropertySet(newPlayerColor).AddCard(targetCard)
		targetPlayer.GetPropertySet(newTargetColor).AddCard(playerCard)

		return ActionResult{
			Success: true,
			Message: fmt.Sprintf("Swapped %s for %s", playerCard.GetName(), targetCard.GetName()),
		}
	}

	return ActionResult{Success: false, Error: "unknown pending action type"}
}

// ========== Helper functions ==========

func (e *Executor) drawCard() Card {
	if len(e.state.Deck) == 0 {
		e.state.Deck = e.state.DiscardPile
		e.state.DiscardPile = nil
		e.shuffleDeck()
	}

	if len(e.state.Deck) == 0 {
		return nil
	}

	card := e.state.Deck[len(e.state.Deck)-1]
	e.state.Deck = e.state.Deck[:len(e.state.Deck)-1]
	return card
}

func (e *Executor) shuffleDeck() {
	for i := len(e.state.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		e.state.Deck[i], e.state.Deck[j] = e.state.Deck[j], e.state.Deck[i]
	}
}

// Helper to parse int from action data
func getIntFromData(data map[string]string, key string, defaultVal int) int {
	if val, ok := data[key]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
