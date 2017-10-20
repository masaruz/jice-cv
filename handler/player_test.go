package handler_test

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"fmt"
	"testing"
	"time"
)

func TestReducer01(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	state.GS.Gambit.Start()
	if !state.GS.Gambit.Check(id1) {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	_, p3 := util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Check ||
		p2.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id2, 30) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Call ||
		p1.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Check ||
		p2.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Call ||
		p3.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Call ||
		p1.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Check ||
		p2.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if state.GS.Gambit.Check(id1) || !state.GS.Gambit.Bet(id1, 40) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Call ||
		p2.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Call ||
		p3.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if state.GS.Gambit.Check(id2) || !state.GS.Gambit.Fold(id2) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Stand {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Call ||
		p3.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if state.GS.Gambit.Check(id3) || !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Stand {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	p1.Print()
	p2.Print()
	p3.Print()
	fmt.Println("now:", time.Now().Unix())
	fmt.Println("end:", state.GS.FinishRoundTime)
}
