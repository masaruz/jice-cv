package api_test

import (
	"999k_engine/api"
	"999k_engine/gambit"
	"999k_engine/handler"
	"999k_engine/state"
	"testing"
)

func TestLoop01(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:  minimumBet,
		BlindsBig:    minimumBet,
		MaxPlayers:   6,
		MaxAFKCount:  5,
		DecisionTime: decisionTime}
	handler.SetGambit(ninek)
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	state.GS.Gambit.Init() // create seats
	// dumb player
	handler.Sit("player1", 2)
	handler.Sit("player2", 3)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1)
	p1 := &state.GS.Players[2]
	p2 := &state.GS.Players[3]
	p3 := &state.GS.Players[5]
	p4 := &state.GS.Players[1]
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if !state.GS.Gambit.Check(p1.ID) {
		t.Error()
	}
	if handler.ExtendPlayerTimeline(p1.ID, 5) {
		t.Error()
	}
	body, err := api.ExtendActionTime(p2.ID)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully extend action time"}` {
		t.Error(data)
	}
	if !handler.ExtendPlayerTimeline(p2.ID, 5) {
		t.Error()
	}
	body, err = api.CashBack(p3.ID)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully cashback"}` {
		t.Error(data)
	}
	if !handler.Stand(p3.ID) {
		t.Error()
	}
	p1.Print()
	p2.Print()
	p3.Print()
	p4.Print()
}
