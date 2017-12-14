package main_test

import (
	"999k_engine/constant"
	"999k_engine/gambit"
	"999k_engine/handler"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLoop01(t *testing.T) {
	decisionTime := int64(1)
	ninek := gambit.NineK{
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		MaxPlayers:      6,
		BuyInMin:        200,
		BuyInMax:        1000,
		BlindsSmall:     10,
		BlindsBig:       10,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	if len(state.Snapshot.Players) != 6 {
		t.Error()
	}
	// dumb player
	handler.Sit("player1", 2)
	if util.CountSitting(state.Snapshot.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.Snapshot.Players) != 2 {
		t.Error()
	}
	handler.Sit("player3", 3)
	if util.CountSitting(state.Snapshot.Players) != 3 {
		t.Error()
	}
	handler.Sit("player4", 1)
	if util.CountSitting(state.Snapshot.Players) != 4 {
		t.Error()
	}
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	// make sure everyone is playing and has 2 cards
	for _, player := range state.Snapshot.Players {
		if player.ID == "" {
			continue
		}
		if player.CardAmount != 2 {
			t.Error()
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
	_, p1 := util.Get(state.Snapshot.Players, "player1")
	_, p2 := util.Get(state.Snapshot.Players, "player2")
	_, p3 := util.Get(state.Snapshot.Players, "player3")
	_, p4 := util.Get(state.Snapshot.Players, "player4")
	newDecisionTime := decisionTime
	if p4.DeadLine-state.Snapshot.StartRoundTime != 4*newDecisionTime ||
		p3.DeadLine-state.Snapshot.StartRoundTime != 2*newDecisionTime ||
		p1.DeadLine-state.Snapshot.StartRoundTime != 1*newDecisionTime ||
		p2.DeadLine-state.Snapshot.StartRoundTime != 3*newDecisionTime {
		t.Error()
	}
	// nothing happend in 2 seconds and assume players act default action
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	// should draw one more card
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
		if player.ID == "" {
			continue
		}
		if player.CardAmount != 3 {
			t.Error()
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
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, "player1")
	// _, p2 := util.Get(state.Snapshot.Players, "player2")
	// _, p3 := util.Get(state.Snapshot.Players, "player3")
	// _, p4 := util.Get(state.Snapshot.Players, "player4")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
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
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	// dumb player
	handler.Sit("player1", 2)
	if util.CountSitting(state.Snapshot.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.Snapshot.Players) != 2 {
		t.Error()
	}
	if state.Snapshot.Gambit.Start() {
		t.Error()
	}
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	if util.CountSitting(state.Snapshot.Players) != 3 {
		t.Error()
	}
	if util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 2 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() || state.Snapshot.Gambit.Start() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	state.Snapshot.Gambit.Start()
	for _, player := range state.Snapshot.Players {
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
	if util.CountSitting(state.Snapshot.Players) != 3 {
		t.Error()
	}
	if util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 3 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, "player1")
	// _, p2 := util.Get(state.Snapshot.Players, "player2")
	// _, p3 := util.Get(state.Snapshot.Players, "player3")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop03(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	// dumb player
	handler.Sit("player1", 2) // dealer
	handler.Sit("player2", 5)
	handler.Sit("player3", 3) // first
	if state.Snapshot.Gambit.Start() {
		t.Error()
	}
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	if state.Snapshot.Gambit.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Check("player2") {
		t.Error()
	}
	// cannot check if already checked
	if state.Snapshot.Gambit.Check("player2") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Check("player1") || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	if state.Snapshot.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Check("player2") || state.Snapshot.Gambit.Check("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Check("player1") {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, "player1")
	// _, p2 := util.Get(state.Snapshot.Players, "player2")
	// _, p3 := util.Get(state.Snapshot.Players, "player3")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime.Unix())
}

func TestLoop04(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2) // dealer
	handler.Sit("player2", 4)
	handler.Sit("player3", 3) // first
	handler.Sit("player4", 5)
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() || util.SumPots(state.Snapshot.PlayerPots) != 40 {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	if !state.Snapshot.Gambit.Check("player3") {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, "player1")
	_, p2 := util.Get(state.Snapshot.Players, "player2")
	_, p3 := util.Get(state.Snapshot.Players, "player3")
	_, p4 := util.Get(state.Snapshot.Players, "player4")
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Bet("player2", 15) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, "player1")
	_, p2 = util.Get(state.Snapshot.Players, "player2")
	_, p3 = util.Get(state.Snapshot.Players, "player3")
	_, p4 = util.Get(state.Snapshot.Players, "player4")
	if util.SumBet(p2) != 25 || p2.Action.Name != constant.Bet ||
		p1.Default.Name != constant.Fold ||
		p3.Default.Name != constant.Fold ||
		p4.Default.Name != constant.Fold {
		t.Error("p2 bets != 25, p2 action name != bet, p1,p3,p4 default != fold")
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Fold("player4") {
		t.Error()
	}
	_, p4 = util.Get(state.Snapshot.Players, "player4")
	if p4.Action.Name != constant.Fold {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if state.Snapshot.Gambit.Check("player1") || !state.Snapshot.Gambit.Call("player3") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.Snapshot.Gambit.Check("player2") {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Bet("player3", 30) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, "player1")
	// _, p2 = util.Get(state.Snapshot.Players, "player2")
	// _, p3 = util.Get(state.Snapshot.Players, "player3")
	// _, p4 = util.Get(state.Snapshot.Players, "player4")
	// p3.Print()
	// p2.Print()
	// p4.Print()
	// p1.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop05(t *testing.T) {
	decisionTime := int64(3)
	delay := int64(0)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
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
	if !state.Snapshot.Gambit.Start() || util.SumPots(state.Snapshot.PlayerPots) != 40 {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Bet("player1", 20) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Call("player2") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+4))
	if !state.Snapshot.Gambit.Bet("player4", 30) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet("player2", 40) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if state.Snapshot.Gambit.Call("player3") {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold("player4") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Call("player1") {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay+3))
	if state.Snapshot.Gambit.Check("player1") {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Bet("player2", 10) {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(delay))
	if !state.Snapshot.Gambit.Call("player1") {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, "player1")
	_, p2 := util.Get(state.Snapshot.Players, "player2")
	if util.SumBet(p1) != util.SumBet(p2) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, "player1")
	// _, p2 = util.Get(state.Snapshot.Players, "player2")
	// _, p3 := util.Get(state.Snapshot.Players, "player3")
	// _, p4 := util.Get(state.Snapshot.Players, "player4")
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop06(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	// dumb player
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	if len(state.Snapshot.Visitors) != 1 {
		t.Error()
	}
	handler.Connect(id2)
	handler.Connect(id3)
	if len(state.Snapshot.Visitors) != 3 {
		t.Error()
	}
	handler.Sit(id1, 2) // first
	if len(state.Snapshot.Visitors) != 2 ||
		util.CountSitting(state.Snapshot.Players) != 1 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 ||
		state.Snapshot.Gambit.Start() {
		t.Error()
	}
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	if len(state.Snapshot.Visitors) != 0 ||
		util.CountSitting(state.Snapshot.Players) != 3 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 {
		t.Error()
	}
	handler.Connect(id4)
	if len(state.Snapshot.Visitors) != 1 ||
		util.CountSitting(state.Snapshot.Players) != 3 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 {
		t.Error()
	}
	handler.Sit(id4, 1) // dealer
	if len(state.Snapshot.Visitors) != 0 ||
		util.CountSitting(state.Snapshot.Players) != 4 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 {
		t.Error()
	}
	handler.Stand(id4)
	handler.Stand(id3)
	if len(state.Snapshot.Visitors) != 2 ||
		util.CountSitting(state.Snapshot.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 {
		t.Error()
	}
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if len(state.Snapshot.Visitors) != 2 ||
		util.CountSitting(state.Snapshot.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 2 {
		t.Error()
	}
	handler.Sit(id4, 1)
	handler.Sit(id3, 5)
	if len(state.Snapshot.Visitors) != 0 ||
		util.CountSitting(state.Snapshot.Players) != 4 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 2 {
		t.Error()
	}
	if state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !handler.Stand(id3) || len(state.Snapshot.Visitors) != 1 ||
		util.CountSitting(state.Snapshot.Players) != 3 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 3 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(id4) || !state.Snapshot.Gambit.Fold(id1) {
		t.Error()
	}
	if !handler.Sit(id3, 3) {
		t.Error()
	}
	_, player3 := util.Get(state.Snapshot.Players, id3)
	if player3.IsPlaying || len(player3.Cards) > 0 {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	// fmt.Println("========== game start ==========")
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// _, p4 := util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error("player1 cannot check")
	}
	// fmt.Println("========== player1 checked ==========")
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Check(id3) {
		t.Error("player3 cannot check")
	}
	// fmt.Println("========== player3 checked ==========")
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Bet(id2, 20) {
		t.Error("player2 cannot bet")
	}
	// fmt.Println("========== player2 bet ==========")
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Call(id4) {
		t.Error("player4 cannot call")
	}
	// fmt.Println("========== player4 called ==========")
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Fold(id1) {
		t.Error("player1 cannot fold")
	}
	// fmt.Println("========== player1 fold ==========")
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	if !state.Snapshot.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	// fmt.Println("========== player3 fold ==========")
	if state.Snapshot.Gambit.Finish() {
		t.Error("able to finish")
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error("unable to go the next round")
	}
	if !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Bet(id4, 20) ||
		!handler.Stand(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	// for _, player := range state.Snapshot.Players {
	// 	player.Print()
	// }
}

func TestLoop07(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	// dumb player
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 3)
	handler.Sit(id3, 4)
	handler.Sit(id4, 1) // dealer
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	// cannot do action again
	if state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id3, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id4) {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// _, p4 := util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	// fmt.Println()
	if !state.Snapshot.Gambit.Fold(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop08(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	// When connect player does not has any chip until they buy in
	for _, player := range state.Snapshot.Visitors {
		if player.Chips != 0 {
			t.Error()
		}
	}
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 5)
	handler.Sit(id3, 1) // dealer
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	for _, player := range state.Snapshot.Visitors {
		if player.Chips != 190 {
			t.Error("player's chips is not equal 990")
		}
	}
	if state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id1) {
		t.Error("player2 can check or player1 cannot check")
	}
	if !state.Snapshot.Gambit.Bet(id2, 20) {
		fmt.Println(state.Snapshot.MaximumBet)
		t.Error("player2 cannot bet 20")
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	_, p3 := util.Get(state.Snapshot.Players, id3)
	if p2.Chips != 170 || util.SumBet(p2) != 30 || p2.Bets[state.Snapshot.Turn] != 20 {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 51) || !state.Snapshot.Gambit.Bet(id3, 50) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p3.Chips != 140 || p3.Bets[state.Snapshot.Turn] != 50 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id1) {
		t.Error("player1 cannot call")
	}
	if !state.Snapshot.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error("able to finish or unable to go to next round")
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error("player1 cannot check")
	}
	if !state.Snapshot.Gambit.Bet(id3, 50) {
		t.Error("player2 cannot bet 50")
	}
	if !state.Snapshot.Gambit.Fold(id1) {
		t.Error("player1 cannot fold")
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error("able to go to next round or unable to finish")
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p3.Chips != 290 || p1.Chips != 140 || p2.Chips != 170 {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop09(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	handler.Sit(id4, 1) // dealer
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !handler.Stand(id2) {
		t.Error("player2 cannot stand")
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error("player1 cannot check")
	}
	if !state.Snapshot.Gambit.Check(id3) {
		t.Error("player3 cannot check")
	}
	if !handler.Stand(id1) {
		t.Error("player1 cannot check")
	}
	if !state.Snapshot.Gambit.Check(id4) {
		t.Error("player3 cannot check")
	}
	if !state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error("cannot go to next round even everyplayer takes action")
	}
	if len(state.Snapshot.Visitors) != 2 ||
		util.CountSitting(state.Snapshot.Players) != 2 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 2 {
		t.Error("visitor != 1, sitting != 2, playing != 2")
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// _, p4 := util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// _, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop10(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.Snapshot.Gambit.Bet(id1, 20) {
		t.Error("player1 cannot bet")
	}
	if !state.Snapshot.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if !state.Snapshot.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.Snapshot.Gambit.Fold(id2) {
		t.Error("player2 cannot fold")
	}
	if !state.Snapshot.Gambit.Fold(id3) {
		t.Error("player3 cannot fold")
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop11(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	// diff with first bet to the pots (-10)
	state.Snapshot.Players[2].Chips = 40
	state.Snapshot.Players[4].Chips = 10
	state.Snapshot.Players[1].Chips = 20
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	if !state.Snapshot.Gambit.Bet(id1, 20) ||
		state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.AllIn ||
		!state.Snapshot.Gambit.AllIn(id2) {
		t.Error(p2.Actions, p2.Chips)
	}
	if !state.Snapshot.Gambit.Call(id3) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	state.Snapshot.Players[2].Cards = []int{48, 49, 50}
	state.Snapshot.Players[4].Cards = []int{28, 29, 30}
	state.Snapshot.Players[1].Cards = []int{32, 33, 34}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	// TODO Until we can check that player actually has not enough chip
	// if state.Snapshot.Gambit.Start() {
	// 	t.Error()
	// }
	// if util.CountSitting(state.Snapshot.Players) != 1 {
	// 	t.Error()
	// }
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop12(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.Snapshot.Gambit.Bet(id1, 20) {
		t.Error("player1 cannot check")
	}
	if state.Snapshot.Gambit.Check(id2) {
		t.Error("player2 can check")
	}
}

func TestLoop13(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2 := "player1", "player2"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2)
	handler.Sit(id2, 4)
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 2 players game should be start")
	}
	handler.Leave(id2)
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() || state.Snapshot.Gambit.Finish() {
		t.Error("can go to next round or unable to finish game")
	}
	handler.Connect(id2)
	handler.Sit(id2, 3)
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[4]
	state.Snapshot.FinishRoundTime = 0
	// Shipping force player to stand
	// index1, _ := util.Get(state.Snapshot.Players, id1)
	// index2, _ := util.Get(state.Snapshot.Players, id2)
	// state.Snapshot.Players[index2].Chips = 9
	// state.Snapshot.FinishRoundTime = 0
	// if state.Snapshot.Gambit.Start() {
	// 	t.Error()
	// }
	// handler.Sit(id2, 3)
	p1.Chips = float64(state.Snapshot.Gambit.GetSettings().BlindsSmall)
	p2.Chips = float64(state.Snapshot.Gambit.GetSettings().BlindsSmall)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 0 {
		t.Error()
	}
	// state.Snapshot.FinishRoundTime = 0
	// state.Snapshot.Gambit.Start()
	// if util.CountSitting(state.Snapshot.Players) != 1 || state.Snapshot.IsGameStart {
	// 	t.Error()
	// }
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// p1.Print()
	// p2.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop14(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        1000,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error("has 3 players game should be start")
	}
	if !state.Snapshot.Gambit.Bet(id1, 20) { // 30 (20+10)
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) { // 30 (20+10)
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id3, 60) { // 70 (60+10)
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id1, 100) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(id3) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	_, p3 := util.Get(state.Snapshot.Players, id3)
	if p1.Chips != 1100 || p2.Chips != 970 || p3.Chips != 930 {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id2, 30) {
		t.Error()
	}
	if state.Snapshot.Gambit.AllIn(id3) || !state.Snapshot.Gambit.Bet(id3, 50) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)

	if p1.Default.Name != constant.Fold || p2.Default.Name != constant.Fold {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id1, 90) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id2, 59) || state.Snapshot.Gambit.Bet(id2, 171) || !state.Snapshot.Gambit.Bet(id2, 170) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 321) || !state.Snapshot.Gambit.Bet(id3, 320) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id2, 770) || !state.Snapshot.Gambit.AllIn(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(id3) {
		t.Error()
	}
	// exceed player's chips
	if state.Snapshot.Gambit.Bet(id1, 721) || state.Snapshot.Gambit.Bet(id1, 589) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id1, 590) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop15(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        1000,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	handler.Sit(id4, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	i1, p1 := util.Get(state.Snapshot.Players, id1)
	i2, p2 := util.Get(state.Snapshot.Players, id2)
	i3, p3 := util.Get(state.Snapshot.Players, id3)
	i4, p4 := util.Get(state.Snapshot.Players, id4)
	// 2100 + 700 + 1300 + 600 = 4700
	state.Snapshot.Players[i1].Chips = 2090
	state.Snapshot.Players[i2].Chips = 690
	state.Snapshot.Players[i3].Chips = 1290
	state.Snapshot.Players[i4].Chips = 590
	if state.Snapshot.Gambit.AllIn(id1) || state.Snapshot.Gambit.Raise(id1, 41) || !state.Snapshot.Gambit.Raise(id1, 40) {
		t.Error()
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 80 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id2, 80) {
		t.Error()
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 160 {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 79) || state.Snapshot.Gambit.Raise(id3, 161) || !state.Snapshot.Gambit.Raise(id3, 160) {
		t.Error(state.Snapshot.MaximumBet)
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 320 {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id4, 159) || state.Snapshot.Gambit.Raise(id4, 321) || !state.Snapshot.Gambit.Raise(id4, 320) {
		t.Error()
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 640 {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id1, 279) || state.Snapshot.Gambit.Raise(id1, 601) || !state.Snapshot.Gambit.Raise(id1, 600) {
		t.Error(state.Snapshot.MinimumBet)
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 1240 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id3) {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(id4) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id1, 9) || !state.Snapshot.Gambit.Bet(id1, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id3, 40) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	state.Snapshot.Players[i1].Cards = []int{40, 41, 42}
	state.Snapshot.Players[i2].Cards = []int{44, 45, 46}
	state.Snapshot.Players[i3].Cards = []int{36, 37, 38}
	state.Snapshot.Players[i4].Cards = []int{48, 49, 50}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	_, p4 = util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	// fmt.Println("pots:", state.Snapshot.PlayerPots)
	// fmt.Println()
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	_, p4 = util.Get(state.Snapshot.Players, id4)
	if p1.Chips != 1430 || p1.IsWinner ||
		p2.Chips != 260 || !p2.IsWinner ||
		p3.Chips != 610 || p3.IsWinner ||
		p4.Chips != 2400 || !p4.IsWinner {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
	// fmt.Println("pots:", state.Snapshot.PlayerPots)
}

func TestLoop16(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	handler.Sit(id4, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	i1, p1 := util.Get(state.Snapshot.Players, id1)
	i2, p2 := util.Get(state.Snapshot.Players, id2)
	i3, p3 := util.Get(state.Snapshot.Players, id3)
	i4, p4 := util.Get(state.Snapshot.Players, id4)
	state.Snapshot.Players[i1].Chips = 390
	state.Snapshot.Players[i2].Chips = 690
	state.Snapshot.Players[i3].Chips = 1290
	state.Snapshot.Players[i4].Chips = 590
	if !state.Snapshot.Gambit.Bet(id1, 20) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 790) || !state.Snapshot.Gambit.Bet(id3, 70) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id4, 69) || !state.Snapshot.Gambit.Bet(id4, 140) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 790) || !state.Snapshot.Gambit.Bet(id3, 450) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id4) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	if p1.Actions[0].Name != constant.Fold || p1.Actions[1].Name != constant.AllIn {
		t.Error()
	}
	if !state.Snapshot.Gambit.AllIn(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Call(id2) || state.Snapshot.Gambit.Bet(id2, 9) || state.Snapshot.Gambit.Raise(id2, 10) || !state.Snapshot.Gambit.Bet(id2, 10) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id3, 9) || !state.Snapshot.Gambit.Bet(id3, 20) {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id4, 19) || !state.Snapshot.Gambit.AllIn(id4) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id3) {
		t.Error()
	}
	state.Snapshot.Players[i1].Cards = []int{40, 41, 42}
	state.Snapshot.Players[i2].Cards = []int{44, 45, 46}
	state.Snapshot.Players[i3].Cards = []int{36, 37, 38}
	state.Snapshot.Players[i4].Cards = []int{48, 49, 50}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	_, p4 = util.Get(state.Snapshot.Players, id4)
	if p1.Chips != 0 || p2.Chips != 100 || p3.Chips != 700 || p4.Chips != 2200 {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("fin:", state.Snapshot.FinishRoundTime)
}

func TestLoop17(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3, id4 := "player1", "player2", "player3", "player4"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.Connect(id4)
	handler.StartTable()
	if state.Snapshot.Gambit.Start() {
		t.Error("not enough players to start the game")
	}
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 5)
	handler.Sit(id4, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	handler.Leave(id1)
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	handler.Leave(id2)
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	handler.Leave(id3)
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// _, p4 := util.Get(state.Snapshot.Players, id4)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
}

func TestLoop18(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	state.Snapshot.Gambit.Start()
	if state.Snapshot.MinimumBet != minimumBet || state.Snapshot.MaximumBet != util.SumPots(state.Snapshot.PlayerPots) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	_, p3 := util.Get(state.Snapshot.Players, id3)
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
	if !state.Snapshot.Gambit.Bet(id2, 20) {
		t.Error()
	}
	if state.Snapshot.MinimumBet != minimumBet+10 {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
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
		p3.Actions[2].Hints[0].Value != 40 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id3) {
		t.Error()
	}
	if state.Snapshot.MinimumBet != minimumBet+10 || state.Snapshot.MaximumBet != 70 {
		t.Error(state.Snapshot.MaximumBet)
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
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
	if state.Snapshot.Gambit.Check(id1) || !state.Snapshot.Gambit.Bet(id1, 40) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
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
		p2.Actions[2].Hints[0].Value != 60 {
		t.Error()
	}
	if state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Fold(id2) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
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
	if state.Snapshot.Gambit.Check(id3) || !state.Snapshot.Gambit.Call(id3) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.MinimumBet != minimumBet || state.Snapshot.MaximumBet != 130 {
		t.Error(state.Snapshot.MaximumBet)
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
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
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop19(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	state.Snapshot.Gambit.Start()
	if state.Snapshot.MinimumBet != minimumBet || state.Snapshot.MaximumBet != util.SumPots(state.Snapshot.PlayerPots) {
		t.Error()
	}
	i1, _ := util.Get(state.Snapshot.Players, id1)
	i2, _ := util.Get(state.Snapshot.Players, id2)
	i3, _ := util.Get(state.Snapshot.Players, id3)
	state.Snapshot.Players[i1].Chips = 390
	state.Snapshot.Players[i2].Chips = 690
	state.Snapshot.Players[i3].Chips = 1290
	if !state.Snapshot.Gambit.Bet(id1, 30) {
		t.Error()
	}
	if state.Snapshot.Gambit.Raise(id2, 30) || !state.Snapshot.Gambit.Raise(id2, 60) {
		t.Error()
	}
	if state.Snapshot.Gambit.Raise(id3, 121) || !state.Snapshot.Gambit.Raise(id3, 120) {
		t.Error()
	}
	if state.Snapshot.Gambit.Raise(id1, 211) || !state.Snapshot.Gambit.Raise(id1, 210) {
		t.Error()
	}
	if state.Snapshot.Gambit.Raise(id2, 420) || state.Snapshot.Gambit.Raise(id2, 179) || !state.Snapshot.Gambit.Raise(id2, 330) {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop20(t *testing.T) {
	t.Run("score=10000002", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		score := scores[0] + scores[1]
		if score != 10000002 || kind != constant.ThreeOfAKind {
			t.Error()
		}
	})
	// t.Run("b>a", func(t *testing.T) {
	// 	as, akind := gambit.NineK{}.Evaluate(model.Cards{
	// 		49, 50, 51,
	// 	})
	// 	bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
	// 		4, 5, 6,
	// 	})
	// 	a := as[0] + as[1]
	// 	b := bs[0] + bs[1]
	// 	if a != 10000051 || b != 10000052 || a > b ||
	// 		akind != constant.ThreeOfAKind ||
	// 		bkind != constant.ThreeOfAKind {
	// 		t.Error()
	// 	}
	// })
	t.Run("a>b", func(t *testing.T) {
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
			4, 17, 20,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a != 10000002 || bs[0] != 6 || a < b ||
			akind != constant.ThreeOfAKind ||
			bkind != constant.Nothing {
			t.Error()
		}
	})
}

func TestLoop21(t *testing.T) {
	t.Run("8c,9c,10c is correct", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			24, 28, 32,
		})
		score := scores[0] + scores[1]
		if score != 1000032 || kind != constant.StraightFlush {
			t.Error()
		}
	})
	t.Run("5c,7c,8c is wrong", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			12, 20, 24,
		})
		if scores[0] != 1000 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("Kc,Ad,2d is wrong", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			44, 49, 1,
		})
		if scores[0] != 3 || kind != constant.Nothing {
			t.Error()
		}
	})
	t.Run("Ac,2c,3c is wrong", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			48, 0, 4,
		})
		if scores[0] != 1000 || scores[1] != 48 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("2c,3d,4h is wrong", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("2c,3c,4c < 6c,7c,8c", func(t *testing.T) {
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
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
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
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
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
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
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("Jc,Qs,Js", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		score := scores[0] + scores[1]
		if score != 100043 ||
			kind != constant.Royal {
			t.Error()
		}
	})
}

func TestLoop23(t *testing.T) {
	t.Run("Qc,Jd,Js is correct", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			40, 37, 39,
		})
		score := scores[0] + scores[1]
		if score != 100040 || kind != constant.Royal {
			t.Error()
		}
	})
	t.Run("Jc,Jd,Js is not correct because it is three of a kind", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			36, 37, 38,
		})
		score := scores[0] + scores[1]
		if score == 100038 || kind != constant.ThreeOfAKind {
			t.Error()
		}
	})
	t.Run("Jc,Qs,Js < Kc,Qc,Jh", func(t *testing.T) {
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Error()
		}
	})
	t.Run("Kh,Qs,Js > Kc,Qc,Jh", func(t *testing.T) {
		as, akind := gambit.NineK{}.Evaluate(model.Cards{
			46, 43, 39,
		})
		bs, bkind := gambit.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a < b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Error()
		}
	})
}

func TestLoop24(t *testing.T) {
	t.Run("Qc,10d,1s is nothing but has bonus", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			40, 33, 51,
		})
		if scores[0] != 1 ||
			scores[1] != 51 ||
			kind != constant.Nothing {
			t.Error()
		}
	})
	t.Run("Jd,Qd,Ah is nothing", func(t *testing.T) {
		scores, kind := gambit.NineK{}.Evaluate(model.Cards{
			37, 41, 50,
		})
		if scores[0] != 1 ||
			scores[1] != 50 ||
			kind != constant.Nothing {
			t.Error()
		}
	})
}

func TestLoop25(t *testing.T) {
	t.Run("6,2,9 hearts must win 10s,2s,5d", func(t *testing.T) {
		a, akind := gambit.NineK{}.Evaluate(model.Cards{
			18, 2, 30,
		})
		b, bkind := gambit.NineK{}.Evaluate(model.Cards{
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

func TestLoop26(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        1000,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	state.Snapshot.Gambit.Start()
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	_, p3 := util.Get(state.Snapshot.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p1.Actions[2].Hints[0].Value != 10 || p1.Actions[2].Hints[1].Value != 30 {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Check ||
		p2.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p2.Actions[2].Hints[0].Value != 10 || p2.Actions[2].Hints[1].Value != 30 {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p3.Actions[2].Hints[0].Value != 10 || p2.Actions[2].Hints[1].Value != 30 {
		t.Error()
	}
	if state.Snapshot.Gambit.Bet(id1, 9) || state.Snapshot.Gambit.Bet(id1, 31) || !state.Snapshot.Gambit.Bet(id1, 30) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p1.Actions[2].Hints[0].Value != 10 || p1.Actions[2].Hints[1].Value != 30 {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Call ||
		p2.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p2.Actions[2].Hints[0].Value != 60 || p2.Actions[2].Hints[1].Value != 60 {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Call ||
		p3.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p3.Actions[2].Hints[0].Value != 60 || p3.Actions[2].Hints[1].Value != 60 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id2, 60) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id3, 120) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Call ||
		p1.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p1.Actions[2].Hints[0].Value != 210 || p1.Actions[2].Hints[1].Value != 210 {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Call ||
		p2.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p2.Actions[2].Hints[0].Value != 180 || p2.Actions[2].Hints[1].Value != 180 {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	// _, p1 = util.Get(state.Snapshot.Players, id1)
	// _, p2 = util.Get(state.Snapshot.Players, id2)
	// _, p3 = util.Get(state.Snapshot.Players, id3)
	// fmt.Println(p1.ID, p1.Bets, p1.Actions)
	// fmt.Println(p2.ID, p2.Bets, p2.Actions)
	// fmt.Println(p3.ID, p3.Bets, p3.Actions)
	// fmt.Println()
	if !state.Snapshot.Gambit.Call(id1) {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p1.Actions[2].Hints[0].Value != 120 || p1.Actions[2].Hints[1].Value != 240 {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Fold ||
		p2.Actions[1].Name != constant.Call ||
		p2.Actions[2].Name != constant.Raise {
		t.Error()
	}
	if p2.Actions[2].Hints[0].Value != 180 || p2.Actions[2].Hints[1].Value != 270 {
		t.Error()
	}
	if p3.Actions[0].Name != constant.Fold ||
		p3.Actions[1].Name != constant.Check ||
		p3.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p3.Actions[2].Hints[0].Value != 120 || p3.Actions[2].Hints[1].Value != 330 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id2) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	_, p3 = util.Get(state.Snapshot.Players, id3)
	if p1.Actions[0].Name != constant.Fold ||
		p1.Actions[1].Name != constant.Check ||
		p1.Actions[2].Name != constant.Bet {
		t.Error()
	}
	if p1.Actions[2].Hints[0].Value != 10 || p1.Actions[2].Hints[1].Value != 390 {
		t.Error()
	}
	// fmt.Println(p1.ID, p1.Bets, p1.Actions[2].Hints[0])
	// fmt.Println(p2.ID, p2.Bets, p2.Actions)
	// fmt.Println(p3.ID, p3.Bets, p3.Actions)
	// fmt.Println(state.Snapshot.MinimumBet)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop27(t *testing.T) {
	decisionTime := int64(1)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	time.Sleep(time.Second * time.Duration(decisionTime))
	if !state.Snapshot.Gambit.Bet(id2, 30) {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	_, p3 := util.Get(state.Snapshot.Players, id3)
	if p1.Action.Name != "" || p1.Default.Name != constant.Fold ||
		p2.Action.Name != constant.Bet || p2.Default.Name != constant.Bet ||
		p3.Action.Name != "" || p3.Default.Name != constant.Fold {
		t.Error()
	}
	if util.SumPots(state.Snapshot.PlayerPots) != 60 {
		t.Error()
	}
}

func TestLoop28(t *testing.T) {
	decisionTime := int64(1)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        1000,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 2) // first
	handler.Sit(id2, 4)
	handler.Sit(id3, 1)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !handler.Stand(id3) || util.SumPots(state.Snapshot.PlayerPots) != 30 ||
		util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) != 2 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id2, 30) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id1, 60) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(id2, 90) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id1) || util.SumPots(state.Snapshot.PlayerPots) != 270 {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// delay after game finish
	if state.Snapshot.FinishRoundTime-time.Now().Unix() != 5 {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	if (p1.Chips != 870 || p2.Chips != 1140) && (p1.Chips != 1140 || p2.Chips != 870) {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop29(t *testing.T) {
	decisionTime := int64(1)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 5) // first
	handler.Sit(id2, 3) // dealer
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	_, p2 := util.Get(state.Snapshot.Players, id2)
	if p2.Type != constant.Dealer || p1.Type != constant.Normal {
		t.Error()
	}
	if !handler.Stand(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	if p2.Type != constant.Dealer {
		t.Error()
	}
	if !p2.IsWinner || p1.IsWinner {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !handler.Sit(id1, 5) || !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	if p1.Type != constant.Dealer || p2.Type != constant.Normal {
		t.Error()
	}
	if !handler.Stand(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	if p2.Type != constant.Normal {
		t.Error()
	}
	if !p2.IsWinner || p1.IsWinner {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !handler.Sit(id1, 5) || !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	if p2.Type != constant.Dealer || p1.Type != constant.Normal {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) || !state.Snapshot.Gambit.Check(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) || !state.Snapshot.Gambit.Check(id2) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	_, p1 = util.Get(state.Snapshot.Players, id1)
	_, p2 = util.Get(state.Snapshot.Players, id2)
	if p1.Type != constant.Dealer || p2.Type != constant.Normal {
		t.Error()
	}
}

func TestLoop30(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 3) // first
	handler.Sit(id2, 5)
	handler.Sit(id3, 1) // dealer
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) || !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id3) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id1) || !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id3) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() || state.Snapshot.Gambit.Start() || state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if state.Snapshot.Gambit.Start() || state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(id2) || !state.Snapshot.Gambit.Check(id3) || !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if state.Snapshot.Gambit.Start() || !state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.Gambit.Start() || state.Snapshot.Gambit.NextRound() || state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	// _, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	// _, p3 := util.Get(state.Snapshot.Players, id3)
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop31(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2, id3 := "player1", "player2", "player3"
	handler.Connect(id1)
	handler.Connect(id2)
	handler.Connect(id3)
	handler.StartTable()
	// dumb player
	handler.Sit(id1, 5) // first
	handler.Sit(id2, 3)
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	state.Snapshot.Players[5].Chips = 30
	state.Snapshot.Players[3].Chips = 120
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id2, 20) {
		t.Error()
	}
	_, p1 := util.Get(state.Snapshot.Players, id1)
	// _, p2 := util.Get(state.Snapshot.Players, id2)
	if len(p1.Actions) > 3 {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}
func TestLoop32(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2 := "player1", "player2"
	handler.Connect(id1)
	handler.Connect(id2)
	// dumb player
	handler.Sit(id1, 5) // first
	handler.Sit(id2, 3)
	p1 := &state.Snapshot.Players[5]
	p2 := &state.Snapshot.Players[3]
	if p1.Actions[0].Name != constant.Stand || p1.Actions[1].Name != constant.StartTable {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Stand || p2.Actions[1].Name != constant.StartTable {
		t.Error()
	}
	handler.Stand(id2)
	if p1.Actions[0].Name != constant.Stand || len(p1.Actions) > 1 {
		t.Error()
	}
	handler.Sit(id2, 3)
	if p1.Actions[0].Name != constant.Stand || p1.Actions[1].Name != constant.StartTable {
		t.Error()
	}
	if p2.Actions[0].Name != constant.Stand || p2.Actions[1].Name != constant.StartTable {
		t.Error()
	}
	handler.Leave(id2)
	if p1.Actions[0].Name != constant.Stand || len(p1.Actions) > 1 {
		t.Error()
	}
	// fmt.Println(p1)
	// fmt.Println(p2)
	// p1.Print()
	// p2.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop33(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	id1, id2 := "player1", "player2"
	handler.Connect(id1)
	handler.Connect(id2)
	// dumb player
	handler.Sit(id1, 5) // first
	handler.Sit(id2, 3)
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	// p1 := &state.Snapshot.Players[5]
	// p2 := &state.Snapshot.Players[3]
	if !state.Snapshot.Gambit.Check(id1) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(id2, 10) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(id1) {
		t.Error()
	}
}

func TestLoop34(t *testing.T) {
	decisionTime := int64(1)
	ninek := gambit.NineK{
		MaxAFKCount:  1,
		MaxPlayers:   6,
		BuyInMin:     1000,
		BuyInMax:     1000,
		BlindsSmall:  10,
		BlindsBig:    10,
		DecisionTime: decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	if len(state.Snapshot.Players) != 6 {
		t.Error()
	}
	// dumb player
	handler.Sit("player1", 2)
	if util.CountSitting(state.Snapshot.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.Snapshot.Players) != 2 {
		t.Error()
	}
	handler.Sit("player3", 3)
	if util.CountSitting(state.Snapshot.Players) != 3 {
		t.Error()
	}
	handler.Sit("player4", 1)
	if util.CountSitting(state.Snapshot.Players) != 4 {
		t.Error()
	}
	handler.StartTable()
	state.Snapshot.Gambit.Start()
	// make sure everyone is playing and has 2 cards
	for _, player := range state.Snapshot.Players {
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
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[5]
	p3 := &state.Snapshot.Players[3]
	p4 := &state.Snapshot.Players[1]
	newDecisionTime := decisionTime
	if p4.DeadLine-state.Snapshot.StartRoundTime != 4*newDecisionTime ||
		p3.DeadLine-state.Snapshot.StartRoundTime != 2*newDecisionTime ||
		p1.DeadLine-state.Snapshot.StartRoundTime != 1*newDecisionTime ||
		p2.DeadLine-state.Snapshot.StartRoundTime != 3*newDecisionTime {
		t.Error()
	}
	// nothing happend in 2 seconds and assume players act default action
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	// should draw one more card
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
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
	time.Sleep(time.Second * time.Duration(state.Snapshot.FinishRoundTime-state.Snapshot.StartRoundTime))
	if state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !handler.Sit("player1", 2) || util.CountSitting(state.Snapshot.Players) != 1 {
		t.Error()
	}
	if !handler.Sit("player2", 5) || util.CountSitting(state.Snapshot.Players) != 2 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
	// fmt.Println(state.Snapshot.AFKCounts, state.Snapshot.DoActions)
}

func TestLoop35(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	cap := 0.50
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		BuyInMin:        1000,
		BuyInMax:        1000,
		DecisionTime:    decisionTime,
		Rake:            5.00,
		Cap:             cap}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	handler.Sit("player1", 2)
	if util.CountSitting(state.Snapshot.Players) != 1 {
		t.Error()
	}
	handler.Sit("player2", 5)
	if util.CountSitting(state.Snapshot.Players) != 2 {
		t.Error()
	}
	handler.Sit("player3", 3)
	if util.CountSitting(state.Snapshot.Players) != 3 {
		t.Error()
	}
	handler.Sit("player4", 1)
	if util.CountSitting(state.Snapshot.Players) != 4 {
		t.Error()
	}
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[5]
	p3 := &state.Snapshot.Players[3]
	p4 := &state.Snapshot.Players[1]
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	for pid, rake := range state.Snapshot.Rakes {
		if pid != "" && rake != 0.5 {
			t.Error()
		}
	}
	if !state.Snapshot.Gambit.Bet(p1.ID, 10) {
		t.Error()
	}
	if state.Snapshot.Rakes[p1.ID] != 1 ||
		state.Snapshot.Rakes[p2.ID] != 0.5 ||
		state.Snapshot.Rakes[p3.ID] != 0.5 ||
		state.Snapshot.Rakes[p4.ID] != 0.5 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(p3.ID, 30) {
		t.Error()
	}
	if state.Snapshot.Rakes[p1.ID] != 1 ||
		state.Snapshot.Rakes[p2.ID] != 0.5 ||
		state.Snapshot.Rakes[p3.ID] != 2 ||
		state.Snapshot.Rakes[p4.ID] != 0.5 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p2.ID) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 0.9090909090909092) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.4545454545454546) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(p4.ID) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 0.9090909090909092) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.4545454545454546) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(p1.ID, 60) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 2.3529411764705883) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.1764705882352942) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.1764705882352942) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.29411764705882354) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p3.ID) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 1.9047619047619047) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 0.9523809523809523) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.9047619047619047) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.23809523809523808) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p2.ID) {
		t.Error()
	}
	if state.Snapshot.Rakes[p1.ID] != 1.6 ||
		state.Snapshot.Rakes[p2.ID] != 1.6 ||
		state.Snapshot.Rakes[p3.ID] != 1.6 ||
		state.Snapshot.Rakes[p4.ID] != 0.2 {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(p1.ID, 40) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 2.0689655172413794) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.3793103448275863) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.3793103448275863) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.1724137931034483) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p3.ID) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.2121212121212122) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.8181818181818183) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.15151515151515152) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p2.ID) {
		t.Error()
	}
	if !util.FloatEquals(state.Snapshot.Rakes[p1.ID], 1.6216216216216217) ||
		!util.FloatEquals(state.Snapshot.Rakes[p2.ID], 1.6216216216216217) ||
		!util.FloatEquals(state.Snapshot.Rakes[p3.ID], 1.6216216216216217) ||
		!util.FloatEquals(state.Snapshot.Rakes[p4.ID], 0.13513513513513514) {
		t.Error()
	}
	if !util.FloatEquals(util.SumRakes(state.Snapshot.Rakes), cap*float64(minimumBet)) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
		if player.IsWinner {
			if player.WinLossAmount != 250 {
				t.Error()
			}
		} else if player.ID != "" {
			if player.ID == "player4" && player.WinLossAmount != -10 {
				t.Error()
			} else if player.ID != "player4" && player.WinLossAmount != -120 {
				t.Error()
			}
		}
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println(state.Snapshot.Rakes)
	// fmt.Println(state.Snapshot.PlayerPots)
}

func TestLoop36(t *testing.T) {
	state.GS = state.GameState{
		TableID:   "default",
		GameIndex: 0,
		Deck:      model.Deck{Cards: model.Cards{0, 1, 2, 3, 4, 5}},
		Players: model.Players{
			model.Player{ID: "player1"},
			model.Player{ID: "player2"},
			model.Player{ID: "player3"},
		},
		Visitors: model.Players{
			model.Player{ID: "player4"},
		},
	}
	if state.Snapshot.TableID != "" {
		t.Error(state.Snapshot)
	}
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	if state.Snapshot.TableID != "default" {
		t.Error()
	}
	state.Snapshot.Players[0].ID = "changed"
	if state.GS.Players[0].ID != "changed" {
		t.Error()
	}
	state.Snapshot.Deck.Cards[2] = 6
	if state.GS.Deck.Cards[2] != 6 {
		t.Error()
	}
	state.Snapshot.Players[0].ID = "player1"
	state.Snapshot.Deck.Cards[2] = 2
	state.Snapshot = state.GameState{
		PlayerTableKeys: make(map[string]string),
		Env:             os.Getenv(constant.Env),
	}
	if state.GS.Players[0].ID != "player1" || state.Snapshot.Players != nil {
		t.Error()
	}
	if state.GS.Deck.Cards[2] != 2 || state.Snapshot.Deck.Cards != nil {
		t.Error()
	}
	state.Snapshot = util.CloneState(state.GS)
	if state.Snapshot.TableID != "default" {
		t.Error()
	}
	state.Snapshot.Players[0].ID = "changed"
	if state.GS.Players[0].ID == "changed" || state.Snapshot.Players[0].ID != "changed" {
		t.Error()
	}
	if state.GS.Deck.Cards[2] != 2 || state.Snapshot.Deck.Cards[2] != 2 {
		t.Error(state.GS.Deck.Cards, state.Snapshot.Deck.Cards)
	}
	state.Snapshot.Deck.Cards[2] = 6
	if state.GS.Deck.Cards[2] == 6 || state.Snapshot.Deck.Cards[2] != 6 {
		t.Error()
	}
}

func TestLoop37(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2)
	handler.Sit("player2", 3)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1)
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[3]
	p3 := &state.Snapshot.Players[5]
	p4 := &state.Snapshot.Players[1]
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) {
		t.Error()
	}
	if handler.ExtendPlayerTimeline(p1.ID) {
		t.Error()
	}
	time.Sleep(time.Second * 1)
	if !handler.ExtendPlayerTimeline(p2.ID) {
		t.Error()
	}
	// gap between is equal how long we sleep
	if p2.StartLine-p1.DeadLine != 1 ||
		// timeline still valid
		p3.StartLine != p2.DeadLine ||
		// other decision time is not changed
		p3.DeadLine-p3.StartLine != 3 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p2.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Bet(p3.ID, 20) {
		t.Error()
	}
	time.Sleep(time.Second * 2)
	if !handler.ExtendPlayerTimeline(p4.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Raise(p4.ID, 40) {
		t.Error()
	}
	time.Sleep(time.Microsecond * 1500)
	if !handler.ExtendPlayerTimeline(p1.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p1.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Call(p2.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Fold(p3.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	// check if new timeline is equal to player amount * decision
	if state.Snapshot.FinishRoundTime-time.Now().Unix() !=
		decisionTime*int64(util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players)) {
		t.Error()
	}
	if p1.Action.Name != "" {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) {
		t.Error()
	}
	if p1.Action.Name != constant.Check {
		t.Error()
	}
	if !handler.ExtendPlayerTimeline(p2.ID) {
		t.Error()
	}
	if p1.Action.Name != constant.Check {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}

func TestLoop38(t *testing.T) {
	/**
	 * Test send stickers
	 */
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2)
	handler.Sit("player2", 3)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1)
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[3]
	p3 := &state.Snapshot.Players[5]
	p4 := &state.Snapshot.Players[1]
	handler.StartTable()
	if state.Snapshot.GameIndex != 0 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) ||
		!state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) ||
		!state.Snapshot.Gambit.Check(p4.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) ||
		!state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) ||
		!state.Snapshot.Gambit.Check(p4.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.GameIndex != 1 {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if state.Snapshot.GameIndex != 2 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) ||
		!state.Snapshot.Gambit.Check(p4.ID) ||
		!state.Snapshot.Gambit.Check(p1.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) ||
		!state.Snapshot.Gambit.Check(p4.ID) ||
		!state.Snapshot.Gambit.Check(p1.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	if state.Snapshot.GameIndex != 2 {
		t.Error()
	}
	handler.SendSticker("xxx", "player1", 2)
	handler.SendSticker("xxx", "player2", 2)
	handler.SendSticker("xxx", "player3", 2)
	if len(*p1.Stickers) != 1 ||
		len(*p2.Stickers) != 1 ||
		len(*p3.Stickers) != 1 {
		t.Error()
	}
	handler.SendSticker("xxx", "player1", 2)
	if len(*p1.Stickers) != 2 ||
		len(*p2.Stickers) != 1 ||
		len(*p3.Stickers) != 1 {
		t.Error()
	}
	time.Sleep(time.Second * 4)
	handler.SendSticker("xxx", "player1", 2)
	handler.SendSticker("xxx", "player3", 2)
	if len(*p1.Stickers) != 1 ||
		len(*p2.Stickers) != 0 ||
		len(*p3.Stickers) != 1 {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
}

func TestLoop39(t *testing.T) {
}

func TestLoop40(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2)
	handler.Sit("player2", 3)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1)
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[3]
	p3 := &state.Snapshot.Players[5]
	p4 := &state.Snapshot.Players[1]
	handler.StartTable()
	if state.Snapshot.GameIndex != 0 {
		t.Error()
	}
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) ||
		!state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) ||
		!state.Snapshot.Gambit.Check(p4.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println(time.Now().Unix())
	// fmt.Println(state.Snapshot.FinishRoundTime)
}

func TestLoop41(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:     minimumBet,
		BlindsBig:       minimumBet,
		BuyInMin:        200,
		BuyInMax:        1000,
		MaxPlayers:      6,
		MaxAFKCount:     5,
		FinishGameDelay: 5,
		DecisionTime:    decisionTime}
	handler.Initiate(ninek)
	state.GS.Gambit.Init() // create seats
	state.Snapshot = state.GS
	state.Snapshot.Duration = 1800
	handler.Connect("player1")
	handler.Connect("player2")
	handler.Connect("player3")
	handler.Connect("player4")
	// dumb player
	handler.Sit("player1", 2)
	handler.Sit("player2", 3)
	handler.Sit("player3", 5)
	handler.Sit("player4", 1)
	p1 := &state.Snapshot.Players[2]
	p2 := &state.Snapshot.Players[3]
	p3 := &state.Snapshot.Players[5]
	p4 := &state.Snapshot.Players[1]
	handler.StartTable()
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	for _, player := range state.Snapshot.Players {
		if player.ID == "" {
			continue
		} else if player.ID != "player4" && player.Type != constant.Normal {
			t.Error()
		} else if player.ID == "player4" && player.Type != constant.Dealer {
			t.Error()
		}
	}
	if !handler.Stand(p4.ID) {
		t.Error()
	}
	if p4.Type != constant.Dealer {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) ||
		!state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p1.ID) ||
		!state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if p1.Type != constant.Dealer {
		t.Error()
	}
	if !handler.Stand(p1.ID) {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.Finish() || !state.Snapshot.Gambit.NextRound() {
		t.Error()
	}
	if !handler.Sit("player1", 2) {
		t.Error()
	}
	if p1.Type != constant.Dealer {
		t.Error()
	}
	if !state.Snapshot.Gambit.Check(p2.ID) ||
		!state.Snapshot.Gambit.Check(p3.ID) {
		t.Error()
	}
	if state.Snapshot.Gambit.NextRound() || !state.Snapshot.Gambit.Finish() {
		t.Error()
	}
	state.Snapshot.FinishRoundTime = 0
	if !state.Snapshot.Gambit.Start() {
		t.Error()
	}
	if p2.Type != constant.Dealer {
		t.Error()
	}
	// p1.Print()
	// p2.Print()
	// p3.Print()
	// p4.Print()
	// fmt.Println("now:", time.Now().Unix())
	// fmt.Println("end:", state.Snapshot.FinishRoundTime)
}
