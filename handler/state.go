package handler

import (
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"time"

	"github.com/googollee/go-socket.io"
)

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
	competitors := createSharedState(state.GS.Players)
	visitors := createSharedState(state.GS.Visitors)
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
				GameState: state.PlayerState{
					Player:       player,
					Competitors:  competitors,
					Visitors:     visitors,
					Pots:         []int{util.SumPots(state.GS.Pots)},
					HighestBet:   util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players),
					Version:      state.GS.Version,
					IsTableStart: state.GS.IsTableStart,
					IsGameStart:  state.GS.IsGameStart}},
			Signature: state.Signature{}})
	return string(data)
}

// filter only attributes are able to be shared
func createSharedState(players model.Players) model.Players {
	others := model.Players{}
	for _, player := range players {
		if state.GS.IsGameStart {
			player.Cards = model.Cards{}
		}
		others = append(others, player)
	}
	return others
}
