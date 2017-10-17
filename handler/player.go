package handler

import (
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
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
		if player.ID != "" {
			state.GS.Players[index].Cards = model.Cards{}
			state.GS.Players[index].Bets = []int{}
			state.GS.Players[index].IsWinner = false
			state.GS.Players[index].IsPlaying = true
		}
	}
	return util.CountSitting(state.GS.Players) >= 2
}

// SetOthersDefaultAction make every has default action
func SetOthersDefaultAction(id string, action string) {
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
		// state.GS.FinishGameTime = ShiftFinishGameTime()
	} else {
		state.GS.Gambit.Finish()
	}
	// state.GS.Save()
}

// Check cards and actioned by player who has turn and has the same bet
func Check(id string) bool {
	index, _ := util.Get(state.GS.Players, id)
	// if !player.IsPlayerTurn {
	// 	return false
	// }
	state.GS.Players[index].Action = model.Action{Name: constant.Check}
	for i := range state.GS.Players {
		if !state.GS.Players[i].IsPlaying {
			continue
		}
		// make default action to everyone
		state.GS.Players[i].Action = model.Action{Name: constant.Check}
	}
	// state.GS.FinishGameTime = ShiftFinishGameTime()
	state.GS.Save()
	return true
}
