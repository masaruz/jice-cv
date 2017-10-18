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

// SetPlayer set updated player to index in gamestate
func SetPlayer(index int, player model.Player) {
	state.GS.Players[index] = player
}

// ActionReducer reduce the actions when player act something
func ActionReducer(event string) model.Actions {
	switch event {
	case constant.StartGame:
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Check},
			model.Action{Name: constant.Bet}}
	case constant.Bet:
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Call},
			model.Action{Name: constant.Raise}}
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
		return model.Actions{}
	case constant.Check:
		return model.Actions{}
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
		state.GS.Players[index].Cards = model.Cards{}
		state.GS.Players[index].Bets = []int{}
		state.GS.Players[index].IsWinner = false
		state.GS.Players[index].IsPlaying = true
		state.GS.Players[index].Actions = ActionReducer(constant.StartGame)
		state.GS.Players[index].Default = model.Action{Name: constant.Check}
		state.GS.Players[index].Action = model.Action{}
	}
	return util.CountSitting(state.GS.Players) >= 2
}

// SetOtherDefaultAction make every has default action
func SetOtherDefaultAction(id string, action model.Action) {
	for index, player := range state.GS.Players {
		if !util.InGame(player) {
			continue
		}
		if id != "" && id != player.ID {
			state.GS.Players[index].Default = action
		} else if id == "" {
			state.GS.Players[index].Default = action
		}
	}
}

// SetOtherActions make every has default action
func SetOtherActions(id string, actions model.Actions) {
	for index, player := range state.GS.Players {
		if !util.InGame(player) {
			continue
		}
		if id != "" && id != player.ID {
			state.GS.Players[index].Actions = actions
		} else if id == "" {
			state.GS.Players[index].Actions = actions
		}
	}
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
			state.GS.Players[index].Actions = ActionReducer(constant.Sit)
		}
	}
	return true
}

// Stand when player need to quit
func Stand(id string) {
	_, caller := util.Get(state.GS.Players, id)
	caller.Action = model.Action{Name: constant.Stand}
	caller.Actions = ActionReducer(constant.Stand)
	state.GS.Players = util.Kick(state.GS.Players, caller.ID)
	// if there are more than 1 player are playing
	if util.CountPlaying(state.GS.Players) > 1 {
		// set to everyone has the same actions
		for index := range state.GS.Players {
			if state.GS.Players[index].ID != "" && state.GS.Players[index].ID != id &&
				!state.GS.Players[index].IsPlaying {
				state.GS.Players[index].Actions = ActionReducer(constant.Sit)
			}
		}
		state.GS.Visitors = util.Add(state.GS.Visitors, caller)
		// state.GS.FinishRoundTime = ShiftFinishRoundTime()
	} else {
		state.GS.Gambit.Finish()
	}
	// state.GS.Save()
}

// Check cards and actioned by player who has turn and has the same bet
func Check(id string) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, player := util.Get(state.GS.Players, id)
	state.GS.Players[index].Action = model.Action{Name: constant.Check}
	diff := time.Now().Sub(player.DeadLine)
	ShiftTimeline(diff)
	return true
}

// Fold cards
func Fold(id string) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, player := util.Get(state.GS.Players, id)
	state.GS.Players[index].Action = model.Action{Name: constant.Fold}
	state.GS.Players[index].Actions = ActionReducer(constant.Fold)
	diff := time.Now().Sub(player.DeadLine)
	ShiftTimeline(diff)
	return true
}

// Bet when previous chips are equally but we want to add more chips to the pots
func Bet(id string, chips int, duration int) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, player := util.Get(state.GS.Players, id)
	// added value to the bet in this turn
	state.GS.Players[index].Bets[state.GS.Turn] += chips
	// broadcast to everyone that I bet
	state.GS.Players[index].Action = model.Action{Name: constant.Bet}
	IncreasePots(chips, 0)
	// others automatic set to fold as default
	SetOtherDefaultAction(id, model.Action{Name: constant.Fold})
	// others need to know what to do next
	SetOtherActions(id, ActionReducer(constant.Bet))
	diff := time.Now().Sub(player.DeadLine)
	ShiftTimeline(diff)
	// duration extend the timeline
	ShiftPlayersToEndOfTimeline(id, duration)
	return true
}

// Call make this player to has same the highest bet
func Call(id string, duration int) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	index, player := util.Get(state.GS.Players, id)
	turn := state.GS.Turn
	chips := util.GetHighestBet(state.GS.Players) - player.Bets[turn]
	state.GS.Players[index].Bets[turn] += chips
	state.GS.Players[index].Action = model.Action{Name: constant.Call}
	OverwriteActionWithDefault(index)
	IncreasePots(chips, 0)
	// others need to know what to do next
	SetOtherActions(id, ActionReducer(constant.Bet))
	diff := time.Now().Sub(player.DeadLine)
	ShiftTimeline(diff)
	return true
}

// InvestToPots added bet to everyone base on turn
func InvestToPots(chips int) {
	// initiate bet value to players
	for index := range state.GS.Players {
		if util.InGame(state.GS.Players[index]) {
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, chips)
			IncreasePots(chips, GetCurrentTurn()) // start with first element in pots
		}
	}
}

// OverwriteActionWithDefault overwritten action with default
func OverwriteActionWithDefault(current int) {
	amount := len(state.GS.Players)
	prev := -1
	round := 0
	for round < amount {
		if current == 0 {
			current = amount
		}
		prev = current - 1
		if util.InGame(state.GS.Players[prev]) && state.GS.Players[prev].Action.Name == "" {
			state.GS.Players[prev].Action = state.GS.Players[prev].Default
		}
		round++
		current--
	}
}
