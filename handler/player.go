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
	}
	return util.CountSitting(state.GS.Players) >= 2
}

// SetDefaultAction make every has default action
func SetDefaultAction(id string, action string) {
	daction := model.Action{Name: action}
	for index, player := range state.GS.Players {
		if !player.IsPlaying {
			continue
		}
		if id != "" && id != player.ID {
			state.GS.Players[index].Action = daction
		} else if id == "" {
			state.GS.Players[index].Action = daction
		}
	}
}

// SetActions make every has default action
func SetActions(id string, actions model.Actions) {
	for index, player := range state.GS.Players {
		if !player.IsPlaying {
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
	_, player := util.Get(state.GS.Players, id)
	if !IsPlayerTurn(id) {
		return false
	}
	SetDefaultAction("", constant.Check)
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
	SetDefaultAction(id, constant.Fold)
	// others need to know what to do next
	SetActions(id, ActionReducer(constant.Bet))
	diff := time.Now().Sub(player.DeadLine)
	ShiftTimeline(diff)
	// duration extend the timeline
	ShiftPlayerTimeline(id, duration)
	return true
}
