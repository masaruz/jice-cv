package game

import (
	"999k_engine/constant"
	"999k_engine/handler"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"sort"
	"time"
)

// NineK is 9K
type NineK struct {
	MaxPlayers   int
	DecisionTime int64
	MinimumBet   int
}

// Payload data accessed by continue
type Payload struct {
	ID    string
	Chips int
}

// Init deck and environment variables
func (game NineK) Init() {
	// set the seats
	handler.CreateSeats(game.MaxPlayers)
	handler.SetMinimumBet(game.MinimumBet)
}

// Start game
func (game NineK) Start() bool {
	if handler.IsTableStart() && !handler.IsGameStart() && handler.MakePlayersReady() {
		handler.StartGame()
		handler.SetMinimumBet(game.MinimumBet)
		// let all players bets to the pots
		handler.InvestToPots(game.MinimumBet)
		// start turn
		handler.IncreaseTurn()
		// start new bets
		handler.InvestToPots(0)
		handler.SetOtherActions("", constant.Check)
		handler.SetDealer()
		handler.BuildDeck()
		handler.Shuffle()
		handler.AssignPlayersCheckOrAllIn()
		handler.CreateTimeLine(game.DecisionTime)
		handler.Deal(2, game.MaxPlayers)
		return true
	}
	return false
}

// NextRound game after round by round
func (game NineK) NextRound() bool {
	if !handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound() &&
		util.CountPlayerNotFold(state.GS.Players) > 1 {
		handler.Deal(1, game.MaxPlayers)
		handler.AssignPlayersCheckOrAllIn()
		handler.CreateTimeLine(game.DecisionTime)
		handler.SetMinimumBet(game.MinimumBet)
		handler.InvestToPots(0)
		handler.IncreaseTurn()
		return true
	}
	handler.OverwriteActionToBehindPlayers()
	return false
}

// Check is doing nothing only shift the timeline
func (game NineK) Check(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	if state.GS.Players[index].Bets[state.GS.Turn] <
		util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		return false
	}
	state.GS.Players[index].Default = model.Action{Name: constant.Check}
	state.GS.Players[index].Action = model.Action{Name: constant.Check}
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	handler.OverwriteActionToBehindPlayers()
	handler.ShortenTimeline(diff)
	return true
}

// Bet is raising bet to the target
func (game NineK) Bet(id string, chips int) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	// not less than minimum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips < state.GS.MinimumBet {
		return false
	}
	// not more than maximum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips > state.GS.MaximumBet {
		return false
	}
	// cannot bet more than player's chips
	if state.GS.Players[index].Chips < chips {
		return false
	}
	// added value to the bet in this turn
	state.GS.Players[index].Chips -= chips
	state.GS.Players[index].Bets[state.GS.Turn] += chips
	// broadcast to everyone that I bet
	state.GS.Players[index].Default = model.Action{Name: constant.Bet}
	state.GS.Players[index].Action = model.Action{Name: constant.Bet}
	state.GS.Players[index].Actions = game.Reducer(constant.Check, id)
	// assign minimum bet
	handler.SetMinimumBet(state.GS.Players[index].Bets[state.GS.Turn])
	handler.SetMaximumBet(util.SumBets(state.GS.Players))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	handler.ShiftPlayersToEndOfTimeline(id, game.DecisionTime)
	return true
}

// Raise is raising bet to the target
func (game NineK) Raise(id string, chips int) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	// not less than minimum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips <= state.GS.MinimumBet {
		return false
	}
	return game.Bet(id, chips)
}

// Call is raising bet to the highest bet
func (game NineK) Call(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	chips := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) -
		state.GS.Players[index].Bets[state.GS.Turn]
	// cannot call more than player's chips
	if state.GS.Players[index].Chips < chips || chips == 0 {
		return false
	}
	state.GS.Players[index].Chips -= chips
	state.GS.Players[index].Bets[state.GS.Turn] += chips
	state.GS.Players[index].Default = model.Action{Name: constant.Call}
	state.GS.Players[index].Action = model.Action{Name: constant.Call}
	state.GS.Players[index].Actions = game.Reducer(constant.Check, id)
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	handler.SetMaximumBet(util.SumBets(state.GS.Players))
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	handler.ShortenTimeline(diff)
	return true
}

// AllIn give all chips
func (game NineK) AllIn(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	chips := state.GS.Players[index].Chips
	// not more than maximum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips > state.GS.MaximumBet {
		return false
	}
	state.GS.Players[index].Bets[state.GS.Turn] += chips
	state.GS.Players[index].Chips = 0
	state.GS.Players[index].Default = model.Action{Name: constant.AllIn}
	state.GS.Players[index].Action = model.Action{Name: constant.AllIn}
	state.GS.Players[index].Actions = game.Reducer(constant.Check, id)
	handler.SetMaximumBet(util.SumBets(state.GS.Players))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	if state.GS.Players[index].Bets[state.GS.Turn] >= util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		handler.SetMinimumBet(state.GS.Players[index].Bets[state.GS.Turn])
		handler.ShiftPlayersToEndOfTimeline(id, game.DecisionTime)
	}
	return true
}

// Fold quit the game but still lost bet
func (game NineK) Fold(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	state.GS.Players[index].Default = model.Action{Name: constant.Fold}
	state.GS.Players[index].Action = model.Action{Name: constant.Fold}
	state.GS.Players[index].Actions = game.Reducer(constant.Fold, id)
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	handler.OverwriteActionToBehindPlayers()
	handler.ShortenTimeline(diff)
	return true
}

// Finish game
func (game NineK) Finish() bool {
	// no others to play with or all players have 3 cards but bet is not equal
	if (util.CountPlayerNotFold(state.GS.Players) <= 1 && handler.IsGameStart()) ||
		// if has 3 cards bet equal
		(handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound()) {
		handler.ForceEndTimeline()
		handler.AssignWinners()
		handler.FlushGame()
		return true
	}
	return false
}

// End game
func (game NineK) End() {}

// Reducer reduce the action and when receive the event
func (game NineK) Reducer(event string, id string) model.Actions {
	switch event {
	case constant.Check:
		_, player := util.Get(state.GS.Players, id)
		// maximum will be player's chips if not enough
		maximum := 0
		if state.GS.MaximumBet > player.Chips {
			maximum = player.Chips
		} else {
			maximum = state.GS.MaximumBet
		}
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Check},
			model.Action{Name: constant.Bet,
				Parameters: model.Parameters{
					model.Parameter{
						Name: "amount", Type: "integer"}},
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: state.GS.MinimumBet},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: maximum}}}}
	case constant.Bet:
		// highest bet in that turn
		highestbet := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players)
		// raise must be highest * 2
		raise := highestbet * 2
		// all sum bets
		pots := util.SumBets(state.GS.Players)
		_, player := util.Get(state.GS.Players, id)
		if highestbet <= player.Bets[state.GS.Turn] {
			return game.Reducer(constant.Check, id)
		}
		// no more than pots
		if raise > pots {
			raise = pots
		}
		if player.Chips < highestbet || player.Chips < raise {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.AllIn,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: player.Chips}}}}
		}
		diff := highestbet - player.Bets[state.GS.Turn]
		// maximum will be player's chips if not enough
		maximum := 0
		if state.GS.MaximumBet > player.Chips {
			maximum = player.Chips
		} else {
			maximum = state.GS.MaximumBet
		}
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Call,
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: diff}}},
			model.Action{Name: constant.Raise,
				Parameters: model.Parameters{
					model.Parameter{
						Name: "amount", Type: "integer"}},
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: raise},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: maximum}}}}

	case constant.Fold:
		return model.Actions{
			model.Action{Name: constant.Stand}}
	default:
		return model.Actions{
			model.Action{Name: constant.Sit}}
	}
}

// Evaluate score of player cards
func (game NineK) Evaluate(values []int) (scores []int, kind string) {
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)
	// all cards are same number
	if threeOfAKind(sorted) {
		return summary(constant.ThreeOfAKind, sorted)
	}
	// all cards are in seqence and same suit
	if straightFlush(sorted) {
		return summary(constant.StraightFlush, sorted)
	}
	// all cards are J or Q or K
	if royal(sorted) {
		return summary(constant.Royal, sorted)
	}
	// all cards are in sequnce
	if straight(sorted) {
		return summary(constant.Straight, sorted)
	}
	// all cards are same color and kind
	if flush(sorted) {
		return summary(constant.Flush, sorted)
	}
	return summary(constant.Nothing, sorted)
}
