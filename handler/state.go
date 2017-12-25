package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/googollee/go-socket.io"
)

// RestoreStateData from env
func RestoreStateData() {
	if s := os.Getenv(constant.Visitors); s != "" {
		visitors := &[]api.Visitor{}
		json.Unmarshal([]byte(s), visitors)
		// Convert to player visitor
		for _, visitor := range *visitors {
			index, _ := util.Get(state.GS.Players, visitor.UserID)
			if index != -1 {
				continue
			}
			state.GS.Visitors = append(state.GS.Visitors,
				model.Player{
					ID:   visitor.UserID,
					Name: visitor.DisplayName,
				})
		}
	}
	if s := os.Getenv(constant.Scoreboard); s != "" {
		scoreboard := &[]model.Scoreboard{}
		json.Unmarshal([]byte(s), scoreboard)
		state.GS.Scoreboard = *scoreboard
	}
	// Assign required parameters
	// Get gameindex from hawkeye who awake this container
	state.GS.GameIndex, _ = strconv.Atoi(os.Getenv(constant.GameIndex))
	state.GS.TableID = os.Getenv(constant.TableID)
	state.GS.GroupID = os.Getenv(constant.GroupID)
	state.GS.StartTableTime, _ = strconv.ParseInt(os.Getenv(constant.StartTime), 10, 64)
	state.GS.Duration, _ = strconv.ParseInt(os.Getenv(constant.Duration), 10, 64)
	// If manager send startime means this table already start
	if state.GS.StartTableTime != 0 {
		state.GS.FinishTableTime = state.GS.StartTableTime + state.GS.Duration
		state.GS.IsTableStart = true
	}
}

// ConvertStringToRequestStruct convert string to struct (state.Req)
// ConvertStringToRequestStruct return as value of pointer
func ConvertStringToRequestStruct(msg string) (*state.Req, error) {
	res := state.Req{}
	err := json.Unmarshal([]byte(msg), &res)
	return &res, err // go function can return multiple values
}

// BroadcastGameState send to everyone but will return caller state in string
func BroadcastGameState(so socketio.Socket, event string, owner string) model.Player {
	player := model.Player{}
	competitor := broadcast(so, state.GS.Players, event, owner)
	visitor := broadcast(so, state.GS.Visitors, event, owner)
	// check if player is competitor or visitor
	if competitor.ID != "" {
		player = competitor
	} else if visitor.ID != "" {
		player = visitor
	}
	return player
}

//  broadcast to everyone in array except caller
func broadcast(so socketio.Socket, players model.Players, event string, owner string) model.Player {
	playerstate := model.Player{}
	for _, player := range players {
		if player.ID == "" {
			continue
		}
		// if is other
		if owner != player.ID {
			so.BroadcastTo(player.ID, constant.PushState, CreateResponse(player.ID, event))
		} else {
			playerstate = player
		}
	}
	return playerstate
}

// CreateResponse what each player should see
func CreateResponse(id string, event string) string {
	competitors := CreateSharedState(state.GS.Players)
	_, player := util.Get(state.GS.Players, id)
	actions := model.Actions{}
	if _, c := util.Get(state.GS.Players, id); c.ID != "" {
		actions = c.Actions
	} else if _, v := util.Get(state.GS.Visitors, id); v.ID != "" {
		actions = v.Actions
	}
	// for record latest state
	if event != "" {
		state.GS.Event = event
	}
	// map to playerstate
	data, _ := json.Marshal(
		state.Resp{
			Header: state.Header{Token: "player_token"},
			Payload: state.RespPayload{
				EventName:       state.GS.Event,
				Actions:         actions,
				CurrentTime:     time.Now().Unix(),
				StartRoundTime:  state.GS.StartRoundTime,
				FinishRoundTime: state.GS.FinishRoundTime,
				IsTableExpired:  state.GS.IsTableExpired,
				GameIndex:       state.GS.GameIndex,
				FinishGameDelay: state.GS.Gambit.GetSettings().FinishGameDelay,
				Scoreboard:      state.GS.Scoreboard,
				GameState: state.PlayerState{
					Player:       player,
					Competitors:  competitors,
					Visitors:     state.GS.Visitors,
					Pots:         []int{util.SumPots(state.GS.PlayerPots)},
					SummaryPots:  state.GS.Pots,
					HighestBet:   util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players),
					Version:      state.GS.Version,
					IsTableStart: state.GS.IsTableStart,
					IsGameStart:  state.GS.IsGameStart,
				}},
			Signature: state.Signature{}})
	return string(data)
}

// CreateSharedState filter only attributes are able to be shared
func CreateSharedState(players model.Players) model.Players {
	fight := 0
	// Count player who actually fold
	for _, player := range state.GS.Players {
		if player.Action.Name != constant.Fold &&
			len(player.Cards) > 0 {
			fight++
		}
	}
	others := model.Players{}
	if state.GS.IsGameStart || fight <= 1 {
		// Decide that players should see the cards
		for _, player := range state.GS.Players {
			// If during gameplay or everyone is fold their cards
			player.Cards = model.Cards{}
			others = append(others, player)
		}
	} else {
		for _, player := range state.GS.Players {
			// If call but is winner
			if player.Action.Name == constant.Fold ||
				(player.Action.Name == constant.Call && !player.IsWinner) {
				player.Cards = model.Cards{}
			}
			others = append(others, player)
		}
	}
	return others
}
