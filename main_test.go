package main_test

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"cardgame/model"
	"fmt"
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
		fmt.Println(util.CountPlayerNotFoldAndNotAllIn(handler.GetPlayerState()) <= 1, handler.IsGameStart(),
			handler.IsFullHand(3), handler.BetsEqual(), handler.IsEndRound())
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
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 2 {
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
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 3 {
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
	if state.GS.Gambit.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Check("player2") {
		t.Error()
	}
	// cannot check if already checked
	if state.GS.Gambit.Check("player2") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Check("player1") || !state.GS.Gambit.NextRound() {
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
	if state.GS.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Check("player2") || state.GS.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Check("player1") {
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
	if !state.GS.Gambit.Start() || util.SumBets(state.GS.Players) != 40 {
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
	if !state.GS.Gambit.Check("player3") {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, "player1")
	_, p2 := util.Get(state.GS.Players, "player2")
	_, p3 := util.Get(state.GS.Players, "player3")
	_, p4 := util.Get(state.GS.Players, "player4")
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Bet("player2", 15) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, "player1")
	_, p2 = util.Get(state.GS.Players, "player2")
	_, p3 = util.Get(state.GS.Players, "player3")
	_, p4 = util.Get(state.GS.Players, "player4")
	if util.SumBet(p2) != 25 || p2.Action.Name != constant.Bet ||
		p1.Default.Name != constant.Fold ||
		p3.Default.Name != constant.Fold ||
		p4.Default.Name != constant.Fold {
		t.Error("p2 bets != 25, p2 action name != bet, p1,p3,p4 default != fold")
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Fold("player4") {
		t.Error()
	}
	_, p4 = util.Get(state.GS.Players, "player4")
	if p4.Action.Name != constant.Fold {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if state.GS.Gambit.Check("player1") || !state.GS.Gambit.Call("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.Check("player2") {
		t.Error()
	}
	if !state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Bet("player3", 30) {
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
	if !state.GS.Gambit.Start() || util.SumBets(state.GS.Players) != 40 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Bet("player1", 20) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Call("player2") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if !state.GS.Gambit.Bet("player4", 30) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	if !state.GS.Gambit.Call("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	if !state.GS.Gambit.Bet("player2", 40) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.GS.Gambit.Call("player3") {
		t.Error()
	}
	if !state.GS.Gambit.Fold("player4") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Call("player1") {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+3))
	if state.GS.Gambit.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Bet("player2", 10) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.GS.Gambit.Call("player1") {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, "player1")
	_, p2 := util.Get(state.GS.Players, "player2")
	if util.SumBet(p1) != util.SumBet(p2) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		fmt.Println(util.CountPlayerNotFoldAndNotAllIn(handler.GetPlayerState()) <= 1, handler.IsGameStart(),
			handler.IsFullHand(3), handler.BetsEqual(), handler.IsEndRound())
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
	// fmt.Println("fin:", state.GS.FinishRoundTime)
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
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 ||
		state.GS.Gambit.Start() {
		t.Error()
	}
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 3 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Connect(id4)
	if len(state.GS.Visitors) != 1 ||
		util.CountSitting(state.GS.Players) != 3 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Sit(id4, 1) // dealer
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 4 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 {
		t.Error()
	}
	handler.Stand(id4)
	handler.Stand(id3)
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 {
		t.Error()
	}
	handler.StartTable()
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 2 {
		t.Error()
	}
	handler.Sit(id4, 1)
	handler.Sit(id3, 5)
	if len(state.GS.Visitors) != 0 ||
		util.CountSitting(state.GS.Players) != 4 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 2 {
		t.Error()
	}
	if state.GS.Gambit.Start() {
		t.Error()
	}
	if !state.GS.Gambit.Check(id2) || !state.GS.Gambit.Check(id1) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if !state.GS.Gambit.Check(id2) || !state.GS.Gambit.Check(id1) {
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
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 3 {
		t.Error()
	}
	if !state.GS.Gambit.Fold(id4) || !state.GS.Gambit.Fold(id1) {
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
	if !state.GS.Gambit.Check(id1) {
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
	if !state.GS.Gambit.Check(id3) {
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
	if !state.GS.Gambit.Bet(id2, 20) {
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
	if !state.GS.Gambit.Call(id4) {
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
	if !state.GS.Gambit.Fold(id1) {
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
	if !state.GS.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	// fmt.Println("========== player3 fold ==========")
	if state.GS.Gambit.Finish() {
		t.Error("able to finish")
	}
	if !state.GS.Gambit.NextRound() {
		t.Error("unable to go the next round")
	}
	if !state.GS.Gambit.Check(id2) || !state.GS.Gambit.Bet(id4, 20) ||
		!handler.Stand(id2) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
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
	if !state.GS.Gambit.Bet(id2, 20) {
		fmt.Println(state.GS.MaximumBet)
		t.Error("player2 cannot bet 20")
	}
	_, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	_, p3 := util.Get(state.GS.Players, id3)
	if p2.Chips != 970 || util.SumBet(p2) != 30 || p2.Bets[state.GS.Turn] != 20 {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 51) || !state.GS.Gambit.Bet(id3, 50) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	if p3.Chips != 940 || p3.Bets[state.GS.Turn] != 50 {
		t.Error()
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
	if p3.Chips != 1090 || p1.Chips != 940 || p2.Chips != 970 {
		t.Error()
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
	if !state.GS.Gambit.Check(id1) {
		t.Error("player1 cannot check")
	}
	if !state.GS.Gambit.Check(id3) {
		t.Error("player3 cannot check")
	}
	if !handler.Stand(id1) {
		t.Error("player1 cannot check")
	}
	if !state.GS.Gambit.Check(id4) {
		t.Error("player3 cannot check")
	}
	if !state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error("cannot go to next round even everyplayer takes action")
	}
	if len(state.GS.Visitors) != 2 ||
		util.CountSitting(state.GS.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 2 {
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
	// diff with first bet to the pots (-10)
	state.GS.Players[2].Chips = 40
	state.GS.Players[4].Chips = 10
	state.GS.Players[1].Chips = 20
	// _, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	// _, p3 := util.Get(state.GS.Players, id3)
	if !state.GS.Gambit.Bet(id1, 20) ||
		state.GS.Gambit.Call(id2) {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.AllIn ||
		!state.GS.Gambit.AllIn(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	// if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
	// 	t.Error()
	// }
	// if !state.GS.Gambit.Finish() {
	// 	fmt.Println(util.CountPlayerNotFold(handler.GetPlayerState()) <= 1, handler.IsGameStart(),
	// 		util.CountPlayerNotFoldAndNotAllIn(handler.GetPlayerState()) <= 1,
	// 		handler.IsFullHand(3), state.GS.Gambit.BetsEqual(), handler.IsEndRound())
	// 	t.Error()
	// }
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
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
	handler.Sit(id1, 2)
	handler.Sit(id2, 4)
	if !state.GS.Gambit.Start() {
		t.Error("has 2 players game should be start")
	}
	handler.Disconnect(id2)
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() || state.GS.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	handler.Connect(id2)
	handler.Sit(id2, 3)
	index1, _ := util.Get(state.GS.Players, id1)
	index2, _ := util.Get(state.GS.Players, id2)
	state.GS.Players[index2].Chips = 9
	if state.GS.Gambit.Start() {
		t.Error()
	}
	handler.Sit(id2, 3)
	state.GS.Players[index1].Chips = 10
	state.GS.Players[index2].Chips = 10
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) != 0 {
		t.Error()
	}
	state.GS.Gambit.Start()
	if util.CountSitting(state.GS.Players) != 1 || state.GS.IsGameStart {
		t.Error()
	}
	// _, p1 := util.Get(state.GS.Players, id1)
	// _, p2 := util.Get(state.GS.Players, id2)
	// p1.Print()
	// p2.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop14(t *testing.T) {
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
	if !state.GS.Gambit.Bet(id1, 20) { // 30 (20+10)
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) { // 30 (20+10)
		t.Error()
	}
	if !state.GS.Gambit.Bet(id3, 60) { // 70 (60+10)
		t.Error()
	}
	if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id1, 100) {
		t.Error()
	}
	if !state.GS.Gambit.Fold(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Fold(id3) {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	_, p3 := util.Get(state.GS.Players, id3)
	if p1.Chips != 1100 || p2.Chips != 970 || p3.Chips != 930 {
		t.Error()
	}
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id2, 30) {
		t.Error()
	}
	if state.GS.Gambit.AllIn(id3) || !state.GS.Gambit.Bet(id3, 50) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)

	if p1.Default.Name != constant.Fold || p2.Default.Name != constant.Fold {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id1, 90) {
		t.Error()
	}
	fmt.Println(state.GS.MinimumBet)
	if state.GS.Gambit.Bet(id2, 59) || state.GS.Gambit.Bet(id2, 171) || !state.GS.Gambit.Bet(id2, 170) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 321) || !state.GS.Gambit.Bet(id3, 320) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id1) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id2, 770) || !state.GS.Gambit.AllIn(id2) {
		t.Error()
	}
	if !state.GS.Gambit.AllIn(id3) {
		t.Error()
	}
	// exceed player's chips
	if state.GS.Gambit.Bet(id1, 721) || state.GS.Gambit.Bet(id1, 589) {
		t.Error()
	}
	if !state.GS.Gambit.Bet(id1, 590) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.GS.Players, id1)
	// _, p2 = util.Get(state.GS.Players, id2)
	// _, p3 = util.Get(state.GS.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop15(t *testing.T) {
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
	handler.Sit(id4, 1)
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	i1, p1 := util.Get(state.GS.Players, id1)
	i2, p2 := util.Get(state.GS.Players, id2)
	i3, p3 := util.Get(state.GS.Players, id3)
	i4, p4 := util.Get(state.GS.Players, id4)
	// 400 + 700 + 1300 + 600 = 3000
	state.GS.Players[i1].Chips = 2090
	state.GS.Players[i2].Chips = 690
	state.GS.Players[i3].Chips = 1290
	state.GS.Players[i4].Chips = 590
	if state.GS.Gambit.AllIn(id1) || state.GS.Gambit.Raise(id1, 41) || !state.GS.Gambit.Raise(id1, 40) {
		t.Error()
	}
	if !state.GS.Gambit.Raise(id2, 80) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 79) || state.GS.Gambit.Raise(id3, 161) || !state.GS.Gambit.Raise(id3, 160) {
		t.Error(state.GS.MaximumBet)
	}
	if state.GS.Gambit.Bet(id4, 159) || state.GS.Gambit.Raise(id4, 321) || !state.GS.Gambit.Raise(id4, 320) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id1, 279) || state.GS.Gambit.Raise(id1, 601) || !state.GS.Gambit.Raise(id1, 600) {
		t.Error(state.GS.MinimumBet)
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	if !state.GS.Gambit.AllIn(id4) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Bet(id1, 9) || !state.GS.Gambit.Bet(id1, 20) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Raise(id3, 40) {
		t.Error()
	}
	if !state.GS.Gambit.Fold(id1) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	state.GS.Players[i1].Cards = []int{48, 49, 50}
	state.GS.Players[i2].Cards = []int{44, 45, 46}
	state.GS.Players[i3].Cards = []int{36, 37, 38}
	state.GS.Players[i4].Cards = []int{40, 41, 42}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	_, p4 = util.Get(state.GS.Players, id4)
	if p1.Chips != 1430 || p2.Chips != 2660 || p3.Chips != 610 || p4.Chips != 0 {
		t.Fail()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop16(t *testing.T) {
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
	handler.Sit(id4, 1)
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	i1, p1 := util.Get(state.GS.Players, id1)
	i2, p2 := util.Get(state.GS.Players, id2)
	i3, p3 := util.Get(state.GS.Players, id3)
	i4, p4 := util.Get(state.GS.Players, id4)
	state.GS.Players[i1].Chips = 390
	state.GS.Players[i2].Chips = 690
	state.GS.Players[i3].Chips = 1290
	state.GS.Players[i4].Chips = 590
	if !state.GS.Gambit.Bet(id1, 20) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 790) || !state.GS.Gambit.Bet(id3, 70) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id4, 69) || !state.GS.Gambit.Bet(id4, 140) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id1) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 790) || !state.GS.Gambit.Bet(id3, 450) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id4) {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	if p1.Actions[0].Name != constant.Fold || p1.Actions[1].Name != constant.AllIn {
		t.Error()
	}
	if !state.GS.Gambit.AllIn(id1) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
		t.Error()
	}
	if state.GS.Gambit.Call(id2) || state.GS.Gambit.Bet(id2, 9) || state.GS.Gambit.Raise(id2, 10) || !state.GS.Gambit.Bet(id2, 10) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id3, 9) || !state.GS.Gambit.Bet(id3, 20) {
		t.Error()
	}
	if state.GS.Gambit.Bet(id4, 19) || !state.GS.Gambit.AllIn(id4) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id2) {
		t.Error()
	}
	if !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	state.GS.Players[i1].Cards = []int{40, 41, 42}
	state.GS.Players[i2].Cards = []int{44, 45, 46}
	state.GS.Players[i3].Cards = []int{36, 37, 38}
	state.GS.Players[i4].Cards = []int{48, 49, 50}
	if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.GS.Players, id1)
	_, p2 = util.Get(state.GS.Players, id2)
	_, p3 = util.Get(state.GS.Players, id3)
	_, p4 = util.Get(state.GS.Players, id4)
	if p1.Chips != 0 || p2.Chips != 100 || p3.Chips != 700 || p4.Chips != 2200 {
		t.Fail()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.GS.FinishRoundTime)
}

func TestLoop17(t *testing.T) {
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
	handler.Sit(id4, 1)
	if !state.GS.Gambit.Start() {
		t.Error()
	}
	handler.Disconnect(id1)
	if state.GS.Gambit.Finish() {
		t.Error()
	}
	handler.Disconnect(id2)
	if state.GS.Gambit.Finish() {
		t.Error()
	}
	handler.Disconnect(id3)
	if !state.GS.Gambit.Finish() {
		t.Error()
	}
	if state.GS.Gambit.Finish() {
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
}

func TestLoop18(t *testing.T) {
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
	if state.GS.MinimumBet != minimumBet || state.GS.MaximumBet != util.SumBets(state.GS.Players) {
		t.Error()
	}
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
	if !state.GS.Gambit.Bet(id2, 20) {
		t.Error()
	}
	if state.GS.MinimumBet != minimumBet+10 {
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
	if p3.Actions[1].Hints[0].Name != "amount" ||
		p3.Actions[1].Hints[0].Type != "integer" ||
		p3.Actions[1].Hints[0].Value != 20 ||
		p3.Actions[2].Parameters[0].Name != "amount" ||
		p3.Actions[2].Parameters[0].Type != "integer" ||
		p3.Actions[2].Hints[0].Name != "amount" ||
		p3.Actions[2].Hints[0].Type != "integer" ||
		p3.Actions[2].Hints[0].Value != 21 {
		t.Error()
	}
	if !state.GS.Gambit.Call(id3) {
		t.Error()
	}
	if state.GS.MinimumBet != minimumBet+10 || state.GS.MaximumBet != 70 {
		t.Error(state.GS.MaximumBet)
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
	if p2.Actions[1].Hints[0].Name != "amount" ||
		p2.Actions[1].Hints[0].Type != "integer" ||
		p2.Actions[1].Hints[0].Value != 20 ||
		p2.Actions[2].Parameters[0].Name != "amount" ||
		p2.Actions[2].Parameters[0].Type != "integer" ||
		p2.Actions[2].Hints[0].Name != "amount" ||
		p2.Actions[2].Hints[0].Type != "integer" ||
		p2.Actions[2].Hints[0].Value != 21 {
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
	if state.GS.MinimumBet != minimumBet || state.GS.MaximumBet != 130 {
		t.Error(state.GS.MaximumBet)
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
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.GS.FinishRoundTime)
}

func TestLoop19(t *testing.T) {
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
	if state.GS.MinimumBet != minimumBet || state.GS.MaximumBet != util.SumBets(state.GS.Players) {
		t.Error()
	}
	i1, _ := util.Get(state.GS.Players, id1)
	i2, _ := util.Get(state.GS.Players, id2)
	i3, _ := util.Get(state.GS.Players, id3)
	state.GS.Players[i1].Chips = 390
	state.GS.Players[i2].Chips = 690
	state.GS.Players[i3].Chips = 1290
	if !state.GS.Gambit.Bet(id1, 30) {
		t.Error()
	}
	if state.GS.Gambit.Raise(id2, 30) || !state.GS.Gambit.Raise(id2, 60) {
		t.Error()
	}
	if state.GS.Gambit.Raise(id3, 121) || !state.GS.Gambit.Raise(id3, 120) {
		t.Error()
	}
	if state.GS.Gambit.Raise(id1, 211) || !state.GS.Gambit.Raise(id1, 210) {
		t.Error()
	}
	if state.GS.Gambit.Raise(id2, 420) || state.GS.Gambit.Raise(id2, 179) || !state.GS.Gambit.Raise(id2, 330) {
		t.Error()
	}
	_, p1 := util.Get(state.GS.Players, id1)
	_, p2 := util.Get(state.GS.Players, id2)
	_, p3 := util.Get(state.GS.Players, id3)
	p1.Print()
	p2.Print()
	p3.Print()
	fmt.Println("now:", time.Now().Unix())
	fmt.Println("end:", state.GS.FinishRoundTime)
}

func TestLoop20(t *testing.T) {
	t.Run("score=10000002", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		score := scores[0] + scores[1]
		if score != 10000002 || kind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("b>a", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			49, 50, 51,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 5, 6,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a != 10000051 || b != 10000052 || a > b ||
			akind != constant.ThreeOfAKind ||
			bkind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("a>b", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 17, 20,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a != 10000002 || bs[0] != 6 || a < b ||
			akind != constant.ThreeOfAKind ||
			bkind != constant.Nothing {
			t.Fail()
		}
	})
}

func TestLoop21(t *testing.T) {
	t.Run("8c,9c,10c is correct", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			24, 28, 32,
		})
		score := scores[0] + scores[1]
		if score != 1000032 || kind != constant.StraightFlush {
			t.Error()
		}
	})
	t.Run("5c,7c,8c is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			12, 20, 24,
		})
		if scores[0] != 1000 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("Kc,Ad,2d is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			44, 49, 1,
		})
		if scores[0] != 3 || kind != constant.Nothing {
			t.Error()
		}
	})
	t.Run("Ac,2c,3c is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			48, 0, 4,
		})
		if scores[0] != 1000 || scores[1] != 48 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("2c,3d,4h is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("2c,3c,4c < 6c,7c,8c", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			16, 20, 24,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b || a != 1000008 || b != 1000024 ||
			akind != constant.StraightFlush ||
			bkind != constant.StraightFlush {
			t.Error()
		}
	})
	t.Run("2c,3c,4c > 5c,7c,8c", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			12, 20, 24,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a < b ||
			akind != constant.StraightFlush ||
			bkind != constant.Flush {
			t.Error()
		}
	})
	t.Run("2c,3c,4c < 2d,3d,4d", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			1, 5, 9,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b ||
			akind != constant.StraightFlush ||
			bkind != constant.StraightFlush {
			t.Error()
		}
	})
}

func TestLoop22(t *testing.T) {
	t.Run("2c,3d,4h is collect", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("Jc,Qs,Js", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		score := scores[0] + scores[1]
		if score != 100043 ||
			kind != constant.Royal {
			t.Fail()
		}
	})
}

func TestLoop23(t *testing.T) {
	t.Run("Qc,Jd,Js is correct", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			40, 37, 39,
		})
		score := scores[0] + scores[1]
		if score != 100040 || kind != constant.Royal {
			t.Fail()
		}
	})
	t.Run("Jc,Jd,Js is not correct because it is three of a kind", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			36, 37, 38,
		})
		score := scores[0] + scores[1]
		if score == 100038 || kind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("Jc,Qs,Js < Kc,Qc,Jh", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Fail()
		}
	})
	t.Run("Kh,Qs,Js > Kc,Qc,Jh", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			46, 43, 39,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a < b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Fail()
		}
	})
}

func TestLoop24(t *testing.T) {
	t.Run("Qc,10d,1s is nothing but has bonus", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			40, 33, 51,
		})
		if scores[0] != 1 ||
			scores[1] != 51 ||
			kind != constant.Nothing {
			t.Fail()
		}
	})
	t.Run("Jd,Qd,Ah is nothing", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			37, 41, 50,
		})
		if scores[0] != 1 ||
			scores[1] != 50 ||
			kind != constant.Nothing {
			t.Fail()
		}
	})
}

func TestLoop25(t *testing.T) {
	t.Run("6,2,9 hearts must win 10s,2s,5d", func(t *testing.T) {
		a, akind := game.NineK{}.Evaluate(model.Cards{
			18, 2, 30,
		})
		b, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 12, 16,
		})
		if a[0] != 1000 || akind != constant.Flush {
			t.Error()
		}
		if b[0] != 1000 || bkind != constant.Flush {
			t.Error()
		}
		if a[0]+a[1] < b[0]+b[1] {
			t.Error()
		}
	})
}
