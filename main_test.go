package main_test

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLoop01(t *testing.T) {
	decisionTime := int64(1)
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   10}
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	if len(state.GS.Players) != 6 {
		t.Error()
	}
	// dumb player
	handler.Sit("player1", 2)
	if util.CountSitting(state.GS.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.GS.Players) != 2 {
		t.Error()
	}
	handler.Sit("player3", 3)
	if util.CountSitting(state.GS.Players) != 3 {
		t.Error()
	}
	handler.Sit("player4", 1)
	if util.CountSitting(state.GS.Players) != 4 {
		t.Error()
	}
	handler.StartTable()
	state.GS.Gambit.Start()
	// make sure everyone is playing and has 2 cards
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 2 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	// test timeline
	_, p1 := util.Get(state.GS.Players, "player1")
	_, p2 := util.Get(state.GS.Players, "player2")
	_, p3 := util.Get(state.GS.Players, "player3")
	_, p4 := util.Get(state.GS.Players, "player4")
	newDecisionTime := decisionTime
	if p4.DeadLine-state.GS.StartRoundTime != 4*newDecisionTime ||
		p3.DeadLine-state.GS.StartRoundTime != 2*newDecisionTime ||
		p1.DeadLine-state.GS.StartRoundTime != 1*newDecisionTime ||
		p2.DeadLine-state.GS.StartRoundTime != 3*newDecisionTime {
		t.Error()
	}
	// nothing happend in 2 seconds and assume players act default action
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	// should draw one more card
	if !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 3 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	if state.GS.Gambit.NextRound() {
		t.Error()
	}
	if !state.GS.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, "player1")
	// _, p2 = util.Get(state.GS.Players, "player2")
	// _, p3 = util.Get(state.GS.Players, "player3")
	// _, p4 = util.Get(state.GS.Players, "player4")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
}

func TestLoop02(t *testing.T) {
	decisionTime := int64(1)
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   10}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	// dumb player
	handler.Sit("player1", 2)
	if util.CountSitting(state.GS.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.GS.Players) != 2 {
		t.Error()
	}
	if state.GS.Gambit.Start() {
		t.Error()
	}
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 2 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	handler.Sit("player3", 1)
	if util.CountSitting(state.GS.Players) != 3 {
		t.Error()
	}
	if util.CountPlaying(state.GS.Players) != 2 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	if !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if !player.IsPlaying {
			continue
		}
		if len(player.Cards) != 3 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	if state.GS.Gambit.NextRound() {
		t.Error()
	}
	if !state.GS.Gambit.Finish() {
		t.Error()
	}
	state.GS.Gambit.Start()
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 2 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	if util.CountSitting(state.GS.Players) != 3 {
		t.Error()
	}
	if util.CountPlaying(state.GS.Players) != 3 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	if !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if !player.IsPlaying {
			continue
		}
		if len(player.Cards) != 3 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	time.Sleep(time.Second * time.Duration(state.GS.FinishRoundTime-state.GS.StartRoundTime))
	if state.GS.Gambit.NextRound() {
		t.Error()
	}
	if !state.GS.Gambit.Finish() {
		t.Error()
	}
	// _, p1 := util.Get(state.GS.Players, "player1")
	// _, p2 := util.Get(state.GS.Players, "player2")
	// _, p3 := util.Get(state.GS.Players, "player3")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.GS.FinishRoundTime)
}

func TestLoop03(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	// dumb player
	handler.Sit("player1", 2) // dealer
	handler.Sit("player2", 5)
	handler.Sit("player3", 3) // first
	if state.GS.Gambit.Start() {
		t.Error()
	}
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 2 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	if handler.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player2") {
		t.Error()
	}
	// cannot check if already checked
	if handler.Check("player2") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player1") || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if !player.IsPlaying {
			continue
		}
		if len(player.Cards) != 3 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	time.Sleep(time.Second * time.Duration(delay+decisionTime))
	if handler.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player2") || handler.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player1") {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	// _, p1 := util.Get(state.GS.Players, "player1")
	// _, p2 := util.Get(state.GS.Players, "player2")
	// _, p3 := util.Get(state.GS.Players, "player3")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.GS.FinishRoundTime.Unix())
}

func TestLoop04(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	// dumb player
	handler.Sit("player1", 2) // dealer
	handler.Sit("player2", 4)
	handler.Sit("player3", 3) // first
	handler.Sit("player4", 5)
	handler.StartTable()
	if !state.GS.Gambit.Start() || state.GS.Pots[0] != 40 {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		if len(player.Cards) != 2 {
			t.Error()
		}
		if !player.IsPlaying {
			t.Error()
		}
		if player.Default.Name != constant.Check {
			t.Error()
		}
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Check("player3") {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, "player1")
	_, p2 := util.Get(state.GS.Players, "player2")
	_, p3 := util.Get(state.GS.Players, "player3")
	_, p4 := util.Get(state.GS.Players, "player4")
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Bet("player2", 15, decisionTime) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, "player1")
	_, p2 = util.Get(state.GS.Players, "player2")
	_, p3 = util.Get(state.GS.Players, "player3")
	_, p4 = util.Get(state.GS.Players, "player4")
	if p2.Bets[0] != 25 || p2.Action.Name != constant.Bet ||
		p1.Default.Name != constant.Fold ||
		p3.Default.Name != constant.Fold ||
		p4.Default.Name != constant.Fold {
		t.Error("p2 bets != 25, p2 action name != bet, p1,p3,p4 default != fold")
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Fold("player4") {
		t.Error()
	}
	_, p4 = util.Get(state.GS.Players, "player4")
	if p4.Action.Name != constant.Fold {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if handler.Check("player1") || !handler.Call("player3", 3) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if handler.Check("player2") {
		t.Error()
	}
	if !state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Bet("player3", 30, 3) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, "player1")
	// _, p2 = util.Get(state.GS.Players, "player2")
	// _, p3 = util.Get(state.GS.Players, "player3")
	// _, p4 = util.Get(state.GS.Players, "player4")
	// p3.Print()
	// p2.Print()
	// p4.Print()
	// p1.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.GS.FinishRoundTime.Unix())
}

func TestLoop05(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2) // first
	handler.Sit("player2", 4)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1) // dealer
	handler.StartTable()
	if !state.GS.Gambit.Start() || state.GS.Pots[0] != 40 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Bet("player1", 20, decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Call("player2", decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if !handler.Bet("player4", 30, decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	if !handler.Call("player1", decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	if !handler.Bet("player2", 40, decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if handler.Call("player3", decisionTime) {
		t.Error()
	}
	if !handler.Fold("player4") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Call("player1", decisionTime) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+3))
	if handler.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Bet("player2", 10, decisionTime) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !handler.Call("player1", decisionTime) {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, "player1")
	_, p2 := util.Get(state.GS.Players, "player2")
	if util.SumBet(p1) != util.SumBet(p2) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, "player1")
	// _, p2 = util.Get(state.GS.Players, "player2")
	// _, p3 := util.Get(state.GS.Players, "player3")
	// _, p4 := util.Get(state.GS.Players, "player4")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime.Unix())
}

func TestLoop06(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	// dumb player
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	if len(state.GS.Visitors) != 1 {
		t.Error()
	}
	handler.Connect(id2)
	handler.Connect(id3)
	if len(state.GS.Visitors) != 3 {
		t.Error()
	}
	handler.Sit(id1, 2) // first
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 1 ||
		util.CountPlaying(state.GS.Players) != 0 ||
		state.GS.Gambit.Start() {
		t.Error()
	}
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 3 ||
		util.CountPlaying(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Connect(id4)
	if len(state.GS.Visitors) != 1 ||
		util.CountSitting(state.GS.Players) != 3 ||
		util.CountPlaying(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Sit(id4, 1) // dealer
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 4 ||
		util.CountPlaying(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Stand(id4)
	handler.Stand(id3)
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlaying(state.GS.Players) != 0 {
		t.Error()
	}
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlaying(state.GS.Players) != 2 {
		t.Error()
	}
	handler.Sit(id4, 1)
	handler.Sit(id3, 5)
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 4 ||
		util.CountPlaying(state.GS.Players) != 2 {
		t.Error()
	}
	if state.GS.Gambit.Start() {
		t.Error()
	}
	if !handler.Check(id2) || !handler.Check(id1) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if !handler.Check(id2) || !handler.Check(id1) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if !handler.Stand(id3) || len(state.GS.Visitors) != 1 ||
		util.CountSitting(state.GS.Players) != 3 ||
		util.CountPlaying(state.GS.Players) != 3 {
		t.Error()
	}
	if !handler.Fold(id4) || !handler.Fold(id1) {
		t.Error()
	}
	if !handler.Sit(id3, 3) {
		t.Error()
	}
	_, player3 := util.Get(state.GS.Players, id3)
	if player3.IsPlaying || len(player3.Cards) > 0 {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	_, player2 := util.Get(state.GS.Players, id2)
	if !player2.IsWinner {
		t.Error()
	}
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	// fmt.Println("========== game start ==========")
	// _, p1 := util.Get(state.GS.Players, id1)
	// _, p2 := util.Get(state.GS.Players, id2)
	// _, p3 := util.Get(state.GS.Players, id3)
	// _, p4 := util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Check(id1) {
		t.Error("player1 cannot check")
	}
	// fmt.Println("========== player1 checked ==========")
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Check(id3) {
		t.Error("player3 cannot check")
	}
	// fmt.Println("========== player3 checked ==========")
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Bet(id2, 20, decisionTime) {
		t.Error("player2 cannot bet")
	}
	// fmt.Println("========== player2 bet ==========")
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Call(id4, decisionTime) {
		t.Error("player4 cannot call")
	}
	// fmt.Println("========== player4 called ==========")
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Fold(id1) {
		t.Error("player1 cannot fold")
	}
	// fmt.Println("========== player1 fold ==========")
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	if !handler.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	// fmt.Println("========== player3 fold ==========")
	if state.GS.Gambit.Finish() {
		t.Error("able to finish")
	}
	if !state.GS.Gambit.NextRound() {
		t.Error("unable to go the next round")
	}
	if !handler.Check(id2) || !handler.Bet(id4, 20, decisionTime) ||
		!handler.Stand(id2) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	_, player4 := util.Get(state.GS.Players, id4)
	if !player4.IsWinner {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	// for _, player := range state.GS.Players {
	// 	player.Print()
	// }
}

func TestLoop07(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	// dumb player
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	handler.Sit(id1, 2)    // first
	handler.Sit(id2, 3)
	handler.Sit(id3, 4)
	handler.Sit(id4, 1) // dealer
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if !state.GS.Gambit.Check(id1) {
		t.Error()
	}
	// cannot do action again
	if state.GS.Gambit.Check(id1) {
		t.Error()
	}
	if !state.GS.Gambit.Check(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id3, 20) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id4) {
		t.Error()
	}
	// _, p1 := util.Get(state.GS.Players, id1)
	// _, p2 := util.Get(state.GS.Players, id2)
	// _, p3 := util.Get(state.GS.Players, id3)
	// _, p4 := util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	// fmt.Println()
	if !state.GS.Gambit.Fold(id1) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop08(t *testing.T) {
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
	for _, player := range state.GS.Visitors {
		if player.Chips != 1000 {
			t.Error("player's chips is not equal 1000")
		}
	}
	handler.StartTable()
	if state.GS.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 5)
	handler.Sit(id3, 1) // dealer
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	for _, player := range state.GS.Visitors {
		if player.Chips != 990 {
			t.Error("player's chips is not equal 990")
		}
	}
	if state.GS.Gambit.Check(id2) || !state.GS.Gambit.Check(id1) {
		t.Error("player2 can check or player1 cannot check")
	}
	if !state.GS.Gambit.Bet(id2, 30) {
		t.Error("player2 cannot bet 30")
	}
	_, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	_, p3 := util.Get(state.GS.Players, id3)
	if p2.Chips != 960 || p2.Bets[state.GS.Turn] != 40 {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 1500) {
		t.Error("player3 bet more than 1000")
	}
	if !state.GS.Gambit.Bet(id3, 900) {
		t.Error("player3 cannot bet 900")
	}
	// _, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p3.Chips != 90 || p3.Bets[state.GS.Turn] != 910 {
		t.Error("player3 chips != 90 or player2 bets[0] != 900")
	}
	if !state.GS.Gambit.Call(id1) {
		t.Error("player1 cannot call")
	}
	if !state.GS.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error("able to finish or unable to go to next round")
	}
	if !state.GS.Gambit.Check(id1) {
		t.Error("player1 cannot check")
	}
	if !state.GS.Gambit.Bet(id3, 50) {
		t.Error("player2 cannot bet 50")
	}
	if !state.GS.Gambit.Fold(id1) {
		t.Error("player1 cannot fold")
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error("able to go to next round or unable to finish")
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p3.Chips != 1950 || p1.Chips != 90 || p2.Chips != 960 {
		t.Error("player3 chips != 1950, player1 chips != 90, player2 chips != 960")
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop09(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.StartTable()
	if state.GS.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	handler.Sit(id4, 1) // dealer
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !handler.Stand(id2) {
		t.Error("player2 cannot stand")
	}
	if !handler.Check(id1) {
		t.Error("player1 cannot check")
	}
	if !handler.Check(id3) {
		t.Error("player3 cannot check")
	}
	if !handler.Stand(id1) {
		t.Error("player1 cannot check")
	}
	if !handler.Check(id4) {
		t.Error("player3 cannot check")
	}
	if !state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error("cannot go to next round even everyplayer takes action")
	}
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlaying(state.GS.Players) != 2 {
		t.Error("visitor != 1, sitting != 2, playing != 2")
	}
	// _, p1 := util.Get(state.GS.Players, id1)
	// _, p2 := util.Get(state.GS.Players, id2)
	// _, p3 := util.Get(state.GS.Players, id3)
	// _, p4 := util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// _, p4 = util.Get(state.GS.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop10(t *testing.T) {
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
	if state.GS.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.GS.Gambit.Bet(id1, 20) {
		t.Error("player1 cannot bet")
	}
	if !state.GS.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if !state.GS.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.GS.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if !state.GS.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	// _, p1 := util.Get(state.GS.Players, id1)
	// _, p2 := util.Get(state.GS.Players, id2)
	// _, p3 := util.Get(state.GS.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop11(t *testing.T) {
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
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	for index := range state.GS.Players {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		state.GS.Players[index].Chips = r1.Intn(1000)
	}
	if !handler.AllIn(id1, decisionTime) ||
		!handler.AllIn(id2, decisionTime) ||
		!handler.AllIn(id3, decisionTime) {
		t.Error()
	}
	for _, player := range state.GS.Players {
		if player.ID != "" {
			player.Print()
		}
	}
	fmt.Println("Pots", state.GS.Pots)
}

func TestLoop12(t *testing.T) {
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
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.GS.Gambit.Bet(id1, 20) {
		t.Error("player1 cannot check")
	}
	if state.GS.Gambit.Check(id2) {
		t.Error("player2 can check")
	}
}

func TestLoop13(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   minimumBet}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init() // create seats
	id1, id2 := "player1", "player2"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	if !state.GS.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	handler.Disconnect(id1)
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() || state.GS.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
}
