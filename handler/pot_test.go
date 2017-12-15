package handler_test

import (
	"999k_engine/gambit"
	"999k_engine/handler"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"testing"
)

func TestLoop01(t *testing.T) {
	ui := &state.GameState{
		Pots: make([]model.Pot, 5),
	}
	for i := range ui.Pots {
		ui.Pots[i].Players = make(map[string]bool)
	}
	handler.CalculatePot(ui, "a", 20)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 20 || ui.Pots[0].Ratio != 20 {
		t.Error(ui.Pots[0])
	}
	if handler.SumPots(ui) != 20 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "b", 20)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 40 || ui.Pots[0].Ratio != 20 {
		t.Error(ui.Pots[0])
	}
	if handler.SumPots(ui) != 40 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "c", 25)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 60 || ui.Pots[0].Ratio != 20 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 5 || ui.Pots[1].Ratio != 25 {
		t.Error(ui.Pots[1])
	}
	if handler.SumPots(ui) != 65 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "d", 5)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 20 || ui.Pots[0].Ratio != 5 {
		ui.Pots.Print()
		t.Error()
	}
	if ui.Pots[1].Value != 45 || ui.Pots[1].Ratio != 20 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 5 || ui.Pots[2].Ratio != 25 {
		t.Error(ui.Pots[2])
	}
	if handler.SumPots(ui) != 70 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "e", 10)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 25 || ui.Pots[0].Ratio != 5 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 20 || ui.Pots[1].Ratio != 10 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 30 || ui.Pots[2].Ratio != 20 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 5 || ui.Pots[3].Ratio != 25 {
		t.Error(ui.Pots[3])
	}
	if handler.SumPots(ui) != 80 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "f", 20)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 30 || ui.Pots[0].Ratio != 5 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 25 || ui.Pots[1].Ratio != 10 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 40 || ui.Pots[2].Ratio != 20 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 5 || ui.Pots[3].Ratio != 25 {
		t.Error(ui.Pots[3])
	}
	if handler.SumPots(ui) != 100 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "g", 25)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 35 || ui.Pots[0].Ratio != 5 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 30 || ui.Pots[1].Ratio != 10 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 50 || ui.Pots[2].Ratio != 20 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 10 || ui.Pots[3].Ratio != 25 {
		t.Error(ui.Pots[3])
	}
	if handler.SumPots(ui) != 125 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "h", 5)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 40 || ui.Pots[0].Ratio != 5 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 30 || ui.Pots[1].Ratio != 10 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 50 || ui.Pots[2].Ratio != 20 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 10 || ui.Pots[3].Ratio != 25 {
		t.Error(ui.Pots[3])
	}
	if handler.SumPots(ui) != 130 {
		t.Error(handler.SumPots(ui))
	}
	handler.CalculatePot(ui, "i", 8)
	handler.MergePots(ui)
	if ui.Pots[0].Value != 45 || ui.Pots[0].Ratio != 5 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 21 || ui.Pots[1].Ratio != 8 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 12 || ui.Pots[2].Ratio != 10 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 50 || ui.Pots[3].Ratio != 20 {
		t.Error(ui.Pots[3])
	}
	if ui.Pots[4].Value != 10 || ui.Pots[4].Ratio != 25 {
		t.Error(ui.Pots[4])
	}
	if handler.SumPots(ui) != 138 {
		t.Error(handler.SumPots(ui))
	}
	delete(ui.Pots[0].Players, "d")
	delete(ui.Pots[0].Players, "h")
	handler.MergePots(ui)
	if ui.Pots[0].Value != 66 || ui.Pots[0].Ratio != 8 {
		t.Error(ui.Pots[0])
	}
	if ui.Pots[1].Value != 12 || ui.Pots[1].Ratio != 10 {
		t.Error(ui.Pots[1])
	}
	if ui.Pots[2].Value != 50 || ui.Pots[2].Ratio != 20 {
		t.Error(ui.Pots[2])
	}
	if ui.Pots[3].Value != 10 || ui.Pots[3].Ratio != 25 {
		t.Error(ui.Pots[3])
	}
	if handler.SumPots(ui) != 138 {
		t.Error(handler.SumPots(ui))
	}
	// ui.Pots.Print()
}

func TestLoop02(t *testing.T) {
	decisionTime := int64(1)
	ninek := gambit.NineK{
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		MaxPlayers:      6,
		BuyInMin:        200,
		BuyInMax:        1000,
		BlindsSmall:     10,
		BlindsBig:       10,
		DecisionTime:    decisionTime,
		Rake:            5.00,
		Cap:             0.5}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = util.CloneState(state.GS)
	state.Snapshot.Duration = 1800
	handler.Enter(model.Player{ID: "a"})
	handler.Enter(model.Player{ID: "b"})
	handler.Enter(model.Player{ID: "c"})
	// dumb player
	handler.Sit("a", 2)
	handler.Sit("b", 5)
	handler.Sit("c", 1)
	a := &state.Snapshot.Players[2]
	b := &state.Snapshot.Players[5]
	c := &state.Snapshot.Players[1]
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
		if !player.IsPlaying {
			continue
		}
		// New version of pot
		handler.CalculatePot(&state.Snapshot, player.ID, util.SumBet(player))
		handler.MergePots(&state.Snapshot)
	}
	if state.Snapshot.Pots[0].Value != 30 || state.Snapshot.Pots[0].Ratio != 10 ||
		!state.Snapshot.Pots[0].Players["a"] ||
		!state.Snapshot.Pots[0].Players["b"] ||
		!state.Snapshot.Pots[0].Players["c"] {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(a.ID, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(b.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(c.ID, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(a.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(b.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Pots[0].Value != 150 || state.Snapshot.Pots[0].Ratio != 50 ||
		!state.Snapshot.Pots[0].Players["a"] ||
		!state.Snapshot.Pots[0].Players["b"] ||
		!state.Snapshot.Pots[0].Players["c"] {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(a.ID, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(b.ID, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(c.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(a.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Pots[0].Value != 230 || state.Snapshot.Pots[0].Ratio != 90 {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	a.Chips = 110
	b.Chips = 120
	c.Chips = 130
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if state.Snapshot.Pots[0].Value != 30 || state.Snapshot.Pots[0].Ratio != 10 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(b.ID, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(c.ID, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(a.ID, 80) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(b.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(c.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(a.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Pots[0].Value != 330 || state.Snapshot.Pots[0].Ratio != 110 {
		t.Error()
	}
	if state.Snapshot.Pots[1].Value != 20 || state.Snapshot.Pots[1].Ratio != 120 {
		t.Error()
	}
	if state.Snapshot.Pots[2].Value != 10 || state.Snapshot.Pots[2].Ratio != 130 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Pots[0].Value != 330 || state.Snapshot.Pots[0].Ratio != 110 {
		t.Error()
	}
	if state.Snapshot.Pots[1].Value != 20 || state.Snapshot.Pots[1].Ratio != 120 {
		t.Error()
	}
	if state.Snapshot.Pots[2].Value != 10 || state.Snapshot.Pots[2].Ratio != 130 {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	a.Chips = 110
	b.Chips = 120
	c.Chips = 130
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(c.ID, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(a.ID, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(b.ID, 80) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(c.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(a.ID) {
		t.Error()
	}
	if !handler.Stand(b.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	handler.Sit("b", 5)
	state.Snapshot.FinishRoundTime = 0
	a.Chips = 110
	b.Chips = 120
	c.Chips = 130
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(a.ID, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(b.ID, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(c.ID, 80) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(a.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(b.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(c.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// a.Print()
	// b.Print()
	// c.Print()
	// state.Snapshot.Pots.Print()
	// log.Println(state.Snapshot.Rakes)
}
