package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"log"
	"os"
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
	caller := model.Player{
		ID:      id,
		Action:  model.Action{Name: constant.Stand},
		Actions: Reducer(constant.Connection, id)}
	state.GS.Visitors = util.Add(state.GS.Visitors, caller)
}

// Enter when player actually call to enter to the room
func Enter(player model.Player) bool {
	if player.ID == "" ||
		player.Name == "" ||
		player.Picture == "" {
		return false
	}
	player.Action = model.Action{Name: constant.Stand}
	player.Actions = Reducer(constant.Connection, player.ID)
	state.GS.Visitors = util.Add(state.GS.Visitors, player)
	return true
}

// Leave and Remove user from vistor or player list
func Leave(id string) bool {
	// force them to stand
	if !Stand(id) {
		return false
	}
	// after they stand then remove from visitor
	state.GS.Visitors = util.Remove(state.GS.Visitors, id)
	if os.Getenv("env") != "dev" {
		body, err := api.RemoveAuth(id)
		log.Println("Response from RemoveAuth", string(body), err)
	}
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
			caller.Type = player.Type
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
	// If not in dev, call api
	if os.Getenv("env") != "dev" {
		// Update buy-in cash
		body, err := api.SaveSettlement(id)
		log.Println("Response from SaveSettlement", string(body), err)
		// Save buy-in cash to real player pocket
		body, err = api.CashBack(id)
		log.Println("Response from CashBack", string(body), err)
	}
	// Change state player to visitor
	visitor := model.Player{
		ID:      id,
		Action:  model.Action{Name: constant.Stand},
		Actions: Reducer(constant.Connection, id),
		Name:    caller.Name,
		Picture: caller.Picture,
	}
	state.GS.Players = util.Kick(state.GS.Players, caller.ID)
	state.GS.Visitors = util.Add(state.GS.Visitors, visitor)
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
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
		if state.GS.Players[pi].Stickers == nil {
			continue
		}
		// Clear all stickers which expired
		for si := 0; si < len(*state.GS.Players[pi].Stickers); si++ {
			stickers := *state.GS.Players[pi].Stickers
			// delay for sticker after expired for 2 sec
			if now-stickers[si].FinishTime >= 2 {
				// update sticker
				stickers = append(stickers[:si], stickers[si+1:]...)
				state.GS.Players[pi].Stickers = &stickers
				si-- // reduce index to prevent out of bound
			}
		}
	}
	// append to stickers array
	if index != -1 {
		sticker := model.Sticker{}
		sticker.StartTime = now
		sticker.FinishTime = now + 2
		sticker.ID = stickerid
		sticker.ToTarget = targetslot
		if state.GS.Players[index].Stickers == nil {
			state.GS.Players[index].Stickers = &[]model.Sticker{}
		}
		stickers := *state.GS.Players[index].Stickers
		stickers = append(stickers, sticker)
		state.GS.Players[index].Stickers = &stickers
	}
}

// GetUserIDFromToken make sure this player has valid table key
// Validate this player has been allowed to access this table
func GetUserIDFromToken(tablekey string) string {
	if os.Getenv("env") == "dev" {
		return "default"
	}
	for userid, key := range state.GS.PlayerTableKeys {
		if key == tablekey {
			log.Printf("Found userid [%s] from tablekey [%s]", userid[:4], tablekey[:4])
			return userid
		}
	}
	if len(tablekey) >= 4 {
		log.Printf("Not found userid from tablekey [%s]", tablekey[:4])
	}
	return ""
}
