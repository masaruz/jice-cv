package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"time"
)

// Reducer reduce player common actions
func Reducer(event string, id string) model.Actions {
	switch event {
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
	default:
		return model.Actions{
			model.Action{Name: constant.Sit}}
	}
}

// Connect and move user as a visitor
func Connect(id string) {
	caller := api.SyncPlayer(id)
	caller.Action = model.Action{Name: constant.Stand}
	caller.Actions = Reducer(constant.Connection, id)
	state.GS.Visitors = util.Add(state.GS.Visitors, caller)
}

// Leave and Remove user from vistor or player list
func Leave(id string) bool {
	// force them to stand
	if !Stand(id) {
		return false
	}
	// after they stand then remove from visitor
	state.GS.Visitors = util.Remove(state.GS.Visitors, id)
	// TODO cashback here
	return true
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
			caller.Name = player.Name
			break
		}
	}
	if caller.Slot == -1 {
		return false
	}
	// TODO buyin here
	// remove from visitor
	state.GS.Visitors = util.Remove(state.GS.Visitors, id)
	caller.Action = model.Action{Name: constant.Sit}
	// add to players
	state.GS.Players[caller.Slot] = caller
	state.GS.AFKCounts[caller.Slot] = 0
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
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
	visitor := api.SyncPlayer(id)
	visitor.Action = model.Action{Name: constant.Stand}
	visitor.Actions = Reducer(constant.Connection, id)
	state.GS.Players = util.Kick(state.GS.Players, caller.ID)
	state.GS.Visitors = util.Add(state.GS.Visitors, visitor)
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
	// TODO cashback here
	return true
}

// SetPlayersRake calculate and set rake for players
func SetPlayersRake(rate float64, cap float64) {
	pots := float64(util.SumPots(state.GS.Pots))
	rake := (rate * pots) / 100
	if rake > cap {
		rake = cap
	}
	for _, player := range state.GS.Players {
		percent := float64(util.SumBet(player)) / pots
		state.GS.Rakes[player.ID] = rake * percent
	}
}

// SendSticker added sticker action to gamestate
func SendSticker(stickerid string, senderid string, targetslot int) {
	index, _ := util.Get(state.GS.Players, senderid)
	// create sticker object
	now := time.Now().Unix()
	// clear expire stickers
	for pi := range state.GS.Players {
		// if no stickers continue
		if len(state.GS.Players[pi].Stickers) <= 0 {
			continue
		}
		for si := 0; si < len(state.GS.Players[pi].Stickers); si++ {
			// delay for sticker after expired for 2 sec
			if now-state.GS.Players[pi].Stickers[si].FinishTime >= 2 {
				// update sticker
				state.GS.Players[pi].Stickers =
					append(state.GS.Players[pi].Stickers[:si],
						state.GS.Players[pi].Stickers[si+1:]...)
				si-- // reduce index to prevent out of bound
			}
		}
	}
	sticker := model.Sticker{}
	sticker.StartTime = now
	sticker.FinishTime = now + 2
	sticker.ID = stickerid
	sticker.ToTarget = targetslot
	// append to stickers array
	if index != -1 {
		state.GS.Players[index].Stickers = append(state.GS.Players[index].Stickers, sticker)
	}
}
