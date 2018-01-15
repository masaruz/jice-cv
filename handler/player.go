package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"log"
	"time"
)

// Reducer reduce player common actions
func Reducer(event string, id string) model.Actions {
	switch event {
	case constant.StartTable:
		index, _ := util.Get(state.Snapshot.Players, id)
		if index != -1 {
			return model.Actions{
				model.Action{Name: constant.Stand}}
		}
		return model.Actions{
			model.Action{Name: constant.Sit}}
	case constant.Sit:
		if !state.Snapshot.IsTableStart &&
			(state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 1 ||
				state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 2) {
			return model.Actions{
				model.Action{Name: constant.Stand},
				model.Action{Name: constant.StartTable}}
		}
		return model.Actions{
			model.Action{Name: constant.Stand}}
	case constant.Connection:
		actions := model.Actions{model.Action{Name: constant.Sit}}
		if !state.Snapshot.IsTableStart &&
			(state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 1 ||
				state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 2) {
			return append(actions, model.Action{Name: constant.StartTable})
		}
		return actions
	default:
		index, _ := util.Get(state.Snapshot.Players, id)
		actions := model.Actions{}
		// If player's sitting
		if index != -1 {
			actions = model.Actions{
				model.Action{Name: constant.Stand}}
		} else {
			actions = model.Actions{
				model.Action{Name: constant.Sit}}
		}
		if !state.Snapshot.IsTableStart &&
			(state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 1 ||
				state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 2) {
			return append(actions, model.Action{Name: constant.StartTable})
		}
		return actions
	}
}

// Connect and move user as a visitor
func Connect(id string) {
	caller := model.Player{
		ID:      id,
		Action:  model.Action{Name: constant.Stand},
		Actions: Reducer(constant.Connection, id)}
	state.Snapshot.Visitors = util.Add(state.Snapshot.Visitors, caller)
}

// Enter when player actually call to enter to the room
func Enter(player model.Player) bool {
	index, _ := util.Get(state.Snapshot.Players, player.ID)
	if index == -1 {
		player.Action = model.Action{Name: constant.Stand}
		player.Actions = Reducer(constant.Connection, player.ID)
		state.Snapshot.Visitors = util.Add(state.Snapshot.Visitors, player)
	}
	return true
}

// Leave and Remove user from vistor or player list
func Leave(id string) bool {
	// force them to stand
	Stand(id, true)
	// after they stand then remove from visitor
	state.Snapshot.Visitors = util.Remove(state.Snapshot.Visitors, id)
	if state.Snapshot.Env != "dev" {
		body, err := api.RemoveAuth(id)
		util.Print("Response from RemoveAuth", string(body), err)
		// Update realtime data ex. Visitors
		body, err = api.UpdateRealtimeData()
		util.Print("Response from UpdateRealtimeData", string(body), err)
	}
	return true
}

// Sit for playing the game
func Sit(id string, slot int) *model.Error {
	index, caller := util.Get(state.Snapshot.Visitors, id)
	caller.Slot = -1
	for _, player := range state.Snapshot.Players {
		if player.ID == "" {
			continue
		}
		// If gps is required then check the distance to others
		if state.Snapshot.Gambit.GetSettings().GPSRestrcited &&
			util.Distance(caller, player) <= 50 {
			util.Print(player.ID, "Is nearby someone")
			return &model.Error{Code: NearOtherPlayers}
		}
	}
	// find slot for them
	for _, player := range state.Snapshot.Players {
		if slot == player.Slot && player.ID == "" {
			caller.Slot = player.Slot
			caller.Type = player.Type
			break
		}
	}
	if caller.Slot == -1 {
		return &model.Error{Code: NoAvailableSeat}
	}
	if state.Snapshot.Env != "dev" {
		body, err := api.CashBack(caller.ID)
		util.Print("Response from Cashback", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		// If cashback error
		if resp.Error != (api.Error{}) && resp.Error.StatusCode != 409 {
			return &model.Error{Code: CashbackError}
		}
		// After cashback success set chips to be 0
		caller.Chips = 0
		// Need request to server for buyin
		body, err = api.BuyIn(caller.ID, state.Snapshot.Gambit.GetSettings().BuyInMin)
		util.Print("Response from BuyIn", string(body), err)
		resp = &api.Response{}
		json.Unmarshal(body, resp)
		// BuyIn must be successful
		if resp.Error != (api.Error{}) {
			err := &model.Error{Code: BuyInError}
			if resp.Error.StatusCode == 422 {
				err = &model.Error{Code: ChipIsNotEnough}
			}
			return err
		}
		util.Print("Buy-in success")
	}
	// Assign how much they buy-in
	caller.Chips = float64(state.Snapshot.Gambit.GetSettings().BuyInMin)
	// Update scoreboard
	UpdateBuyInAmount(&caller)
	// remove from visitor
	state.Snapshot.Visitors = util.Remove(state.Snapshot.Visitors, id)
	caller.Action = model.Action{Name: constant.Sit}
	// add to players
	state.Snapshot.Players[caller.Slot] = caller
	state.Snapshot.AFKCounts[caller.Slot] = 0
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
	// Update realtime data ex. Visitors
	if state.Snapshot.Env != "dev" {
		body, err := api.UpdateRealtimeData()
		util.Print("Response from UpdateRealtimeData", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		if resp.Error != (api.Error{}) {
			return &model.Error{Code: UpdateRealtimeError}
		}
	}
	state.Snapshot.AFKCounts[index] = 0
	return nil
}

// Stand when player need to quit
func Stand(id string, force bool) bool {
	index, caller := util.Get(state.Snapshot.Players, id)
	if caller.ID == "" {
		return false
	}
	// If not in dev, call api
	// If this player already buyin
	// Update buy-in cash
	if state.Snapshot.Env != "dev" {
		body, err := api.SaveSettlement(id)
		util.Print("Response from SaveSettlement", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		if resp.Error != (api.Error{}) && resp.Error.StatusCode != 422 && !force {
			return false
		}
		// Save buy-in cash to real player pocket
		body, err = api.CashBack(id)
		util.Print("Response from CashBack", string(body), err)
		resp = &api.Response{}
		json.Unmarshal(body, resp)
		if resp.Error != (api.Error{}) && resp.Error.StatusCode != 409 && !force {
			return false
		}
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
	// Change state player to visitor
	visitor := model.Player{
		ID:              id,
		Action:          model.Action{Name: constant.Stand},
		Actions:         Reducer(constant.Connection, id),
		Name:            caller.Name,
		AvatarSource:    caller.AvatarSource,
		AvatarBuiltinID: caller.AvatarBuiltinID,
		AvatarCustomID:  caller.AvatarCustomID,
		FacebookID:      caller.FacebookID,
		Lat:             caller.Lat,
		Lon:             caller.Lon,
	}
	state.Snapshot.Players = util.Kick(state.Snapshot.Players, caller.ID)
	state.Snapshot.Visitors = util.Add(state.Snapshot.Visitors, visitor)
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
	// Update realtime data ex. Visitors
	if state.Snapshot.Env != "dev" {
		body, err := api.UpdateRealtimeData()
		util.Print("Response from UpdateRealtimeData", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		if resp.Error != (api.Error{}) && !force {
			return false
		}
	}
	state.Snapshot.AFKCounts[index] = 0
	return true
}

// SetPlayersRake calculate and set rake for players
func SetPlayersRake(rate float64, cap float64) {
	pots := float64(util.SumPots(state.Snapshot.PlayerPots))
	rake := (rate * pots) / 100
	if rake > cap {
		rake = cap
	}
	for _, player := range state.Snapshot.Players {
		if player.IsPlaying {
			percent := float64(util.SumBet(player)) / pots
			state.Snapshot.Rakes[player.ID] = rake * percent
		}
	}
}

// SendSticker added sticker action to gamestate
func SendSticker(stickerid string, senderid string, targetslot int) {
	index, _ := util.Get(state.Snapshot.Players, senderid)
	// create sticker object
	now := time.Now().Unix()
	// clear expire stickers
	for pi := range state.Snapshot.Players {
		// if no stickers continue
		if state.Snapshot.Players[pi].Stickers == nil {
			continue
		}
		// Clear all stickers which expired
		for si := 0; si < len(*state.Snapshot.Players[pi].Stickers); si++ {
			stickers := *state.Snapshot.Players[pi].Stickers
			// delay for sticker after expired for 2 sec
			if now-stickers[si].FinishTime >= 2 {
				// update sticker
				stickers = append(stickers[:si], stickers[si+1:]...)
				state.Snapshot.Players[pi].Stickers = &stickers
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
		if state.Snapshot.Players[index].Stickers == nil {
			state.Snapshot.Players[index].Stickers = &[]model.Sticker{}
		}
		stickers := *state.Snapshot.Players[index].Stickers
		stickers = append(stickers, sticker)
		state.Snapshot.Players[index].Stickers = &stickers
	}
}

// GetUserIDFromToken make sure this player has valid table key
// Validate this player has been allowed to access this table
func GetUserIDFromToken(tablekey string) string {
	if state.Snapshot.Env == "dev" {
		return "default"
	}
	for userid, playerTableKey := range state.GS.PlayerTableKeys {
		if playerTableKey.TableKey == tablekey {
			log.Printf("Found userid [%s] from tablekey [%s]", userid[:4], tablekey[:4])
			return userid
		}
	}
	if len(tablekey) >= 4 {
		log.Printf("Not found userid from tablekey [%s]", tablekey[:4])
	}
	return ""
}

// SaveHistory record winloss amount and cards
func SaveHistory() {
	competitors := CreateSharedCardState(state.Snapshot)
	// Convert competitors to competitors history
	histories := []*model.PlayerHistory{}
	for _, comp := range competitors {
		if comp.ID == "" {
			continue
		}
		histories = append(histories, &model.PlayerHistory{
			ID:            comp.ID,
			Name:          comp.Name,
			WinLossAmount: comp.WinLossAmount,
			Cards:         comp.Cards,
			Slot:          comp.Slot,
		})
	}
	for _, comp := range competitors {
		if comp.ID == "" {
			continue
		}
		_, player := util.Get(state.Snapshot.Players, comp.ID)
		history := model.History{
			Player: &model.PlayerHistory{
				ID:            player.ID,
				Name:          player.Name,
				WinLossAmount: player.WinLossAmount,
				Cards:         player.Cards,
				Slot:          player.Slot,
			},
			Competitors: histories,
		}
		if state.Snapshot.History[comp.ID] == nil {
			state.Snapshot.History[comp.ID] = make(map[int]model.History)
		}
		state.Snapshot.History[comp.ID][state.Snapshot.GameIndex] = history
	}
}

// SetPlayerLocation for validation if needed
func SetPlayerLocation(id string, lat float64, lon float64) {
	player := &model.Player{}
	index, _ := util.Get(state.Snapshot.Players, id)
	if index != -1 {
		player = &state.Snapshot.Players[index]
	} else {
		index, _ = util.Get(state.Snapshot.Visitors, id)
		if index == -1 {
			return
		}
		player = &state.Snapshot.Visitors[index]
	}
	player.Lat = lat
	player.Lon = lon
}
