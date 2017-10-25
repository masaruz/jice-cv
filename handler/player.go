package handler

import (
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"time"
)

// GetPlayerState return players state
func GetPlayerState() model.Players {
	return state.GS.Players
}

// ActionReducer reduce the actions when player act something
func ActionReducer(event string, id string) model.Actions {
	switch event {
	case constant.Check:
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
						Name: "amount_max", Type: "integer", Value: state.GS.MaximumBet}}}}
	case constant.Bet:
		highestbet := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players)
		_, player := util.Get(state.GS.Players, id)
		if highestbet <= player.Bets[state.GS.Turn] {
			return ActionReducer(constant.Check, id)
		}
		if player.Chips < highestbet {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.AllIn,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: player.Chips}}}}
		}
		diff := highestbet - player.Bets[state.GS.Turn]
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
						Name: "amount", Type: "integer", Value: diff + 1},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: state.GS.MaximumBet}}}}
	case constant.Sit:
		if util.CountSitting(state.GS.Players) >= 2 && !state.GS.IsTableStart {
			return model.Actions{
				model.Action{Name: constant.Stand},
				model.Action{Name: constant.StartTable}}
		}
		return model.Actions{
			model.Action{Name: constant.Stand}}
	case constant.StartTable:
		return model.Actions{
			model.Action{Name: constant.Stand}}
	case constant.Fold:
		return model.Actions{
			model.Action{Name: constant.Stand}}
	default:
		return model.Actions{
			model.Action{Name: constant.Sit}}
	}
}

// MakePlayersReady make everyone isPlayer = true
func MakePlayersReady() bool {
	for index, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		// force to stand when player has no chips enough
		if player.Chips < state.GS.MinimumBet {
			Stand(player.ID)
			continue
		}
		state.GS.Players[index].Cards = model.Cards{}
		state.GS.Players[index].Bets = []int{}
		state.GS.Players[index].IsPlaying = true
		state.GS.Players[index].IsEarned = false
		state.GS.Players[index].Actions = ActionReducer(constant.Check, state.GS.Players[index].ID)
		state.GS.Players[index].Default = model.Action{Name: constant.Check}
		state.GS.Players[index].Action = model.Action{}
	}
	return util.CountSitting(state.GS.Players) >= 2
}

// SetOtherDefaultAction make every has default action
func SetOtherDefaultAction(id string, action string) {
	daction := model.Action{Name: action}
	for index, player := range state.GS.Players {
		if !util.IsPlayingAndNotFoldAndNotAllIn(player) {
			continue
		}
		if id != "" && id != player.ID {
			_, caller := util.Get(state.GS.Players, id)
			// if caller's bet more than other then overwrite their action
			if caller.Bets[state.GS.Turn] > state.GS.Players[index].Bets[state.GS.Turn] {
				state.GS.Players[index].Default = daction
			}
		} else if id == "" {
			state.GS.Players[index].Default = daction
		}
	}
}

// SetOtherActions make every has default action
func SetOtherActions(id string, action string) {
	for index, player := range state.GS.Players {
		if !util.IsPlayingAndNotFoldAndNotAllIn(player) {
			continue
		}
		if id != "" && id != player.ID {
			state.GS.Players[index].Actions = ActionReducer(action, player.ID)
		} else if id == "" {
			state.GS.Players[index].Actions = ActionReducer(action, player.ID)
		}
	}
}

// Connect and move user as a visitor
func Connect(id string) {
	caller := util.SyncPlayer(id)
	caller.Action = model.Action{Name: constant.Stand}
	caller.Actions = ActionReducer(constant.Connection, id)
	state.GS.Visitors = util.Add(state.GS.Visitors, caller)
}

// Disconnect and Remove user from vistor or player list
func Disconnect(id string) {
	_, caller := util.Get(state.GS.Players, id)
	// if playing need to do something
	if !caller.IsPlaying {
		// TODO shift timeline
	}
	state.GS.Players = util.Kick(state.GS.Players, id)
	state.GS.Visitors = util.Remove(state.GS.Visitors, id)
}

// AutoSit auto find a seat
func AutoSit(id string) bool {
	for _, player := range state.GS.Players {
		if player.ID == "" {
			return Sit(id, player.Slot)
		}
	}
	return false
}

// Sit for playing the game
func Sit(id string, slot int) bool {
	_, caller := util.Get(state.GS.Visitors, id)
	caller.Slot = -1
	// find slot for them
	for _, player := range state.GS.Players {
		if slot == player.Slot && player.ID == "" {
			caller.Slot = player.Slot
			break
		}
	}
	if caller.Slot == -1 {
		return false
	}
	// remove from visitor
	state.GS.Visitors = util.Remove(state.GS.Visitors, id)
	caller.Action = model.Action{Name: constant.Sit}
	// add to players
	state.GS.Players[caller.Slot] = caller
	// if others who are not playing then able to starttable or only stand
	for index, player := range state.GS.Players {
		// not a seat and not playing
		if player.ID != "" && !player.IsPlaying {
			state.GS.Players[index].Actions = ActionReducer(constant.Sit, player.ID)
		}
	}
	return true
}

// Stand when player need to quit
func Stand(id string) bool {
	_, caller := util.Get(state.GS.Players, id)
	if caller.ID == "" {
		return false
	}
	// if playing need to shift timeline
	if caller.IsPlaying {
		if IsPlayerTurn(id) {
			diff := time.Now().Unix() - caller.DeadLine
			ShortenTimeline(diff)
		} else if !util.IsPlayerBehindTheTimeline(caller) {
			diff := caller.StartLine - caller.DeadLine
			ShortenTimelineAfterTarget(id, diff)
		}
		OverwriteActionToBehindPlayers()
	}
	visitor := util.SyncPlayer(id)
	visitor.Action = model.Action{Name: constant.Stand}
	visitor.Actions = ActionReducer(constant.Connection, id)
	state.GS.Players = util.Kick(state.GS.Players, caller.ID)
	state.GS.Visitors = util.Add(state.GS.Visitors, visitor)
	return true
}

// Check cards and actioned by player who has turn and has the same bet
func Check(id string) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, caller := util.Get(state.GS.Players, id)
	if caller.Bets[state.GS.Turn] < util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		return false
	}
	state.GS.Players[index].Default = model.Action{Name: constant.Check}
	state.GS.Players[index].Action = model.Action{Name: constant.Check}
	diff := time.Now().Unix() - caller.DeadLine
	OverwriteActionToBehindPlayers()
	ShortenTimeline(diff)
	return true
}

// Fold cards
func Fold(id string) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, caller := util.Get(state.GS.Players, id)
	state.GS.Players[index].Default = model.Action{Name: constant.Fold}
	state.GS.Players[index].Action = model.Action{Name: constant.Fold}
	state.GS.Players[index].Actions = ActionReducer(constant.Fold, id)
	diff := time.Now().Unix() - caller.DeadLine
	OverwriteActionToBehindPlayers()
	ShortenTimeline(diff)
	return true
}

// AllIn when player has chips less than highest bet
func AllIn(id string, duration int64) bool {
	if !IsPlayerTurn(id) {
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
	state.GS.Players[index].Actions = ActionReducer(constant.Check, id)
	IncreasePots(chips, 0)
	// set action of everyone
	OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - state.GS.Players[index].DeadLine
	ShortenTimeline(diff)
	// duration extend the timeline
	if state.GS.Players[index].Bets[state.GS.Turn] >= util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		state.GS.MinimumBet = state.GS.Players[index].Bets[state.GS.Turn]
		ShiftPlayersToEndOfTimeline(id, duration)
	}
	return true
}

// Raise when previous chips are bet
func Raise(id string, chips int, duration int64) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	// not less than minimum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips <= state.GS.MinimumBet {
		return false
	}
	return Bet(id, chips, duration)
}

// Bet when previous chips are equally but we want to add more chips to the pots
func Bet(id string, chips int, duration int64) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, caller := util.Get(state.GS.Players, id)
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
	state.GS.Players[index].Actions = ActionReducer(constant.Check, id)
	// assign minimum bet
	state.GS.MinimumBet = state.GS.Players[index].Bets[state.GS.Turn]
	IncreasePots(chips, 0)
	// set action of everyone
	OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - caller.DeadLine
	ShortenTimeline(diff)
	// duration extend the timeline
	ShiftPlayersToEndOfTimeline(id, duration)
	return true
}

// Call make this player to has same the highest bet
func Call(id string, duration int64) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, caller := util.Get(state.GS.Players, id)
	chips := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) - caller.Bets[state.GS.Turn]
	// cannot call more than player's chips
	if state.GS.Players[index].Chips < chips || chips == 0 {
		return false
	}
	state.GS.Players[index].Chips -= chips
	state.GS.Players[index].Bets[state.GS.Turn] += chips
	state.GS.Players[index].Default = model.Action{Name: constant.Call}
	state.GS.Players[index].Action = model.Action{Name: constant.Call}
	state.GS.Players[index].Actions = ActionReducer(constant.Check, id)
	// set action of everyone
	OverwriteActionToBehindPlayers()
	IncreasePots(chips, 0)
	// others need to know what to do next
	SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - caller.DeadLine
	ShortenTimeline(diff)
	return true
}

// OverwriteActionToBehindPlayers overwritten action with default
func OverwriteActionToBehindPlayers() {
	for index := range state.GS.Players {
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) &&
			util.IsPlayerBehindTheTimeline(state.GS.Players[index]) {
			state.GS.Players[index].Action = state.GS.Players[index].Default
		}
	}
}

// BurnBet burn bet from player
func BurnBet(id string, burn int) int {
	index, player := util.Get(state.GS.Players, id)
	// if this player cannot pay all of it
	sumbet := util.SumBet(player)
	if burn >= sumbet {
		for i := range state.GS.Players[index].Bets {
			state.GS.Players[index].Bets[i] = 0
		}
		return sumbet
	}
	for i, bet := range state.GS.Players[index].Bets {
		if bet >= burn {
			state.GS.Players[index].Bets[i] -= burn
		} else {
			burn -= bet
			state.GS.Players[index].Bets[i] = 0
		}
	}
	return util.SumBet(player)
}
