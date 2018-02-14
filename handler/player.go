package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"log"
	"math"
	"time"
)

// Reducer reduce player common actions
func Reducer(event string, id string) model.Actions {
	topupAction := GetTopUpHint(id)
	switch event {
	case constant.StartTable:
		index, _ := util.Get(state.Snapshot.Players, id)
		if index != -1 {
			return model.Actions{
				model.Action{Name: constant.Stand}, topupAction}
		}
		return model.Actions{
			model.Action{Name: constant.Sit},
			model.Action{Name: constant.Stand}}
	case constant.Sit:
		if !state.Snapshot.IsTableStart &&
			(state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 1 ||
				state.Snapshot.PlayerTableKeys[id].ClubMemberLevel == 2) {
			return model.Actions{
				model.Action{Name: constant.Stand},
				model.Action{Name: constant.StartTable},
				topupAction}
		}
		return model.Actions{
			model.Action{Name: constant.Stand},
			topupAction}
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
				model.Action{Name: constant.Stand},
				topupAction}
		} else {
			actions = model.Actions{
				model.Action{Name: constant.Sit},
				model.Action{Name: constant.Stand}}
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
	message := &api.GetPlayerMessage{}
	if state.Snapshot.Env != "dev" {
		body, err := api.GetPlayer(player.ID, state.Snapshot.GroupID)
		util.Print("Response from Get A Player", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		// If get player error
		if resp.Error != (api.Error{}) {
			return false
		}
		json.Unmarshal([]byte(resp.Message), message)
	} else {
		// When testing assign buyin max
		message.Player.Chips = float64(state.Snapshot.Gambit.GetSettings().BuyInMax)
	}
	index, _ := util.Get(state.Snapshot.Players, player.ID)
	if index == -1 {
		player.Action = model.Action{Name: constant.Stand}
		player.Actions = Reducer(constant.Connection, player.ID)
		player.TotalChips = message.Player.Chips
		state.Snapshot.Visitors = util.Add(state.Snapshot.Visitors, player)
	} else {
		state.Snapshot.Players[index].TotalChips = message.Player.Chips
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
		go func() {
			body, err := api.RemoveAuth(id)
			util.Print("Response from RemoveAuth", string(body), err)
			// Update realtime data ex. Visitors
			body, err = api.UpdateRealtimeData()
			util.Print("Response from UpdateRealtimeData", string(body), err)
		}()
	}
	return true
}

// Sit for playing the game
func Sit(id string, slot int) *model.Error {
	index, caller := util.Get(state.Snapshot.Visitors, id)
	if state.Snapshot.Gambit.GetSettings().GPSRestrcited &&
		caller.Lat == 0 && caller.Lon == 0 {
		return &model.Error{Code: NearOtherPlayers}
	}
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
		message := &api.GetPlayerMessage{}
		json.Unmarshal([]byte(resp.Message), message)
		caller.TotalChips = message.Player.Chips
		util.Print("Get player success")
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
		go func() {
			body, err := api.UpdateRealtimeData()
			util.Print("Response from UpdateRealtimeData", string(body), err)
		}()
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
	UpdateWinningsAmount(caller.ID, caller.WinLossAmount)
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
		AvatarTimestamp: caller.AvatarTimestamp,
		FacebookID:      caller.FacebookID,
		Lat:             caller.Lat,
		Lon:             caller.Lon,
	}
	state.Snapshot.Players = util.Kick(state.Snapshot.Players, caller.ID)
	state.Snapshot.Visitors = util.Add(state.Snapshot.Visitors, visitor)
	SetOtherActionsWhoAreNotPlaying(constant.Sit)
	// Update realtime data ex. Visitors
	if state.Snapshot.Env != "dev" {
		go func() {
			body, err := api.UpdateRealtimeData()
			util.Print("Response from UpdateRealtimeData", string(body), err)
		}()
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

// SaveTempHistory save history during game
func SaveTempHistory() {
	comps := CreateSharedCardState(state.Snapshot)
	// Convert competitors to competitors history
	for _, comp := range comps {
		if comp.ID == "" || !comp.IsPlaying {
			continue
		}
		_, player := util.Get(state.Snapshot.Players, comp.ID)
		// Skip player themselves and who is not playing
		histories := []model.PlayerHistory{}
		for _, tmp := range comps {
			if tmp.ID == "" ||
				tmp.ID == player.ID ||
				!tmp.IsPlaying {
				continue
			}
			histories = append(histories, model.PlayerHistory{
				ID:            tmp.ID,
				Name:          tmp.Name,
				WinLossAmount: tmp.WinLossAmount,
				Cards:         tmp.Cards,
				Slot:          tmp.Slot,
				CardAmount:    tmp.CardAmount,
			})
		}
		for tmpID, tmp := range state.Snapshot.TempHistory {
			if tmpID == player.ID || state.Snapshot.GameIndex != tmp.GameIndex {
				continue
			}
			has := false
			for _, his := range histories {
				if tmpID == his.ID {
					has = !has
					break
				}
			}
			if !has {
				// When found player who's not in this game anymore hide their cards
				tmp.Player.Cards = model.Cards{}
				histories = append(histories, tmp.Player)
			}
		}
		state.Snapshot.TempHistory[comp.ID] = model.History{
			Player: model.PlayerHistory{
				ID:            player.ID,
				Name:          player.Name,
				WinLossAmount: player.WinLossAmount,
				Cards:         player.Cards,
				Slot:          player.Slot,
				CardAmount:    player.CardAmount,
			},
			Competitors: histories,
			CreateTime:  time.Now().Unix(),
			GameIndex:   state.Snapshot.GameIndex,
		}
	}
	for tmpIndex, tmp := range state.Snapshot.TempHistory {
		if state.Snapshot.GameIndex == tmp.GameIndex {
			for _, comp := range comps {
				for index, his := range tmp.Competitors {
					if comp.ID == his.ID {
						if comp.IsWinner {
							state.Snapshot.TempHistory[tmpIndex].Competitors[index] =
								state.Snapshot.TempHistory[comp.ID].Player
							continue
						}
						state.Snapshot.TempHistory[tmpIndex].Competitors[index].WinLossAmount =
							state.Snapshot.TempHistory[comp.ID].Player.WinLossAmount
						state.Snapshot.TempHistory[tmpIndex].Competitors[index].CardAmount =
							state.Snapshot.TempHistory[comp.ID].Player.CardAmount
					}
				}
			}
		}
	}
}

// SaveHistory record winloss amount and cards
func SaveHistory() {
	// Backup latest history before save to temp
	for playerid, tmp := range state.Snapshot.TempHistory {
		state.Snapshot.History[playerid] = tmp
	}
	SaveTempHistory()
	for tmpID, tmp := range state.Snapshot.TempHistory {
		for hisID, his := range state.Snapshot.History {
			if tmpID == hisID {
				his.Player = tmp.Player
				for _, tmpCom := range tmp.Competitors {
					for hisComI, hisCom := range his.Competitors {
						if tmpCom.ID == hisCom.ID {
							his.Competitors[hisComI] = tmpCom
						}
					}
				}
				state.Snapshot.History[hisID] = his
			}
		}
	}
	if state.Snapshot.Env != "dev" {
		go func() {
			// Need request to server for buyin
			body, err := api.SaveHistories()
			util.Print("Response from Save Histories", string(body), err)
		}()
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

// GetTopUpHint get hint for topup
func GetTopUpHint(id string) model.Action {
	_, player := util.Get(state.Snapshot.Players, id)
	chip := int(math.Floor(player.Chips))
	topupMax := state.Snapshot.Gambit.GetSettings().BuyInMax - chip - int(math.Floor(player.TopUp.Amount))
	// Never be negative
	if topupMax < 0 {
		topupMax = 0
	} else if total := int(math.Floor(player.TotalChips)); topupMax > total {
		topupMax = total
	}
	return model.Action{
		Name: constant.TopUp,
		Hints: model.Hints{
			model.Hint{Name: "amount", Type: "integer", Value: 0},
			model.Hint{Name: "amount_max", Type: "integer", Value: topupMax},
		}}
}

// PrepareTopUp marked player who request for topup and wait for process next turn
func PrepareTopUp(id string, amount float64) bool {
	index, _ := util.Get(state.Snapshot.Players, id)
	if index == -1 || amount == 0 {
		return false
	}
	player := &state.Snapshot.Players[index]
	player.TopUp.IsRequest = true
	player.TopUp.Amount += amount
	if player.IsPlaying {
		player.Actions = state.Snapshot.Gambit.Reducer(player.Action.Name, player.ID)
	} else {
		player.Actions = Reducer(constant.TopUp, player.ID)
	}
	return true
}

// TopUp add amount and be ready to update chips next turn
func TopUp(id string) *model.Error {
	index, _ := util.Get(state.Snapshot.Players, id)
	if index == -1 {
		return &model.Error{Code: PlayerNotFound}
	}
	player := &state.Snapshot.Players[index]
	// If this player does not request yet
	if !player.TopUp.IsRequest || player.TopUp.Amount == 0 {
		return nil
	}
	amount := player.TopUp.Amount
	player.TopUp.Amount = 0
	player.TopUp.IsRequest = false
	if state.Snapshot.Env != "dev" {
		// Need request to server for buyin
		body, err := api.BuyIn(player.ID, int(math.Floor(amount)))
		util.Print("Response from BuyIn", string(body), err)
		resp := &api.Response{}
		json.Unmarshal(body, resp)
		// BuyIn must be successful
		if resp.Error != (api.Error{}) {
			err := &model.Error{Code: BuyInError}
			if resp.Error.StatusCode == 422 {
				err = &model.Error{Code: ChipIsNotEnough}
			}
			return err
		}
		message := &api.GetPlayerMessage{}
		json.Unmarshal([]byte(resp.Message), message)
		player.TotalChips = message.Player.Chips
		util.Print("Get player success")
	}
	player.Chips += amount
	return nil
}

// GetHistories from a user
func GetHistories(id string) error {
	body, err := api.GetHistories(id)
	util.Print("Response from GetHistories", string(body), err)
	resp := &api.Response{}
	json.Unmarshal(body, resp)
	if resp.Error != (api.Error{}) {
		return err
	}
	message := &struct {
		Histories []model.History `json:"histories"`
	}{}
	err = json.Unmarshal([]byte(resp.Message), message)
	state.Snapshot.Histories[id] = message.Histories
	return nil
}
