package game_test

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"testing"
	"time"
)

func TestLoop1(t *testing.T) {
	t.Run("no actions in this game", func(t *testing.T) {
		decisionTime := 1
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   10}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		if len(state.GS.Players) != 6 {
			t.Fail()
		}
		// dumb player
		handler.Sit("player1", 2)
		if util.CountSitting(state.GS.Players) != 1 {
			t.Fail()
		}
		handler.Sit("player2", 5)
		if util.CountSitting(state.GS.Players) != 2 {
			t.Fail()
		}
		handler.Sit("player3", 3)
		if util.CountSitting(state.GS.Players) != 3 {
			t.Fail()
		}
		handler.Sit("player4", 1)
		if util.CountSitting(state.GS.Players) != 4 {
			t.Fail()
		}
		handler.StartTable()
		state.GS.Gambit.Start()
		// make sure everyone is playing and has 2 cards
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if len(player.Cards) != 2 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		// test timeline
		_, p1 := util.Get(state.GS.Players, "player1")
		_, p2 := util.Get(state.GS.Players, "player2")
		_, p3 := util.Get(state.GS.Players, "player3")
		_, p4 := util.Get(state.GS.Players, "player4")
		newDecisionTime := decisionTime
		if p4.DeadLine.Sub(state.GS.StartRoundTime).Seconds() != float64(4*newDecisionTime) ||
			p3.DeadLine.Sub(state.GS.StartRoundTime).Seconds() != float64(2*newDecisionTime) ||
			p1.DeadLine.Sub(state.GS.StartRoundTime).Seconds() != float64(1*newDecisionTime) ||
			p2.DeadLine.Sub(state.GS.StartRoundTime).Seconds() != float64(3*newDecisionTime) {
			t.Fail()
		}
		// nothing happend in 2 seconds and assume players act default action
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		// should draw one more card
		if !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if state.GS.Gambit.Finish() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if len(player.Cards) != 3 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		if state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if !state.GS.Gambit.Finish() {
			t.Fail()
		}
		// _, p1 = util.Get(state.GS.Players, "player1")
		// _, p2 = util.Get(state.GS.Players, "player2")
		// _, p3 = util.Get(state.GS.Players, "player3")
		// _, p4 = util.Get(state.GS.Players, "player4")
		// p1.Print()
		// p2.Print()
		// p3.Print()
		// p4.Print()
	})
}

func TestLoop2(t *testing.T) {
	t.Run("no actions but someone sit during game", func(t *testing.T) {
		decisionTime := 1
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   10}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		// dumb player
		handler.Sit("player1", 2)
		if util.CountSitting(state.GS.Players) != 1 {
			t.Fail()
		}
		handler.Sit("player2", 5)
		if util.CountSitting(state.GS.Players) != 2 {
			t.Fail()
		}
		if state.GS.Gambit.Start() {
			t.Fail()
		}
		handler.StartTable()
		if !state.GS.Gambit.Start() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if len(player.Cards) != 2 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		handler.Sit("player3", 1)
		if util.CountSitting(state.GS.Players) != 3 {
			t.Fail()
		}
		if util.CountPlaying(state.GS.Players) != 2 {
			t.Fail()
		}
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		if !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if state.GS.Gambit.Finish() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if !player.IsPlaying {
				continue
			}
			if len(player.Cards) != 3 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		if state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if !state.GS.Gambit.Finish() {
			t.Fail()
		}
		state.GS.Gambit.Start()
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if len(player.Cards) != 2 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		if util.CountSitting(state.GS.Players) != 3 {
			t.Fail()
		}
		if util.CountPlaying(state.GS.Players) != 3 {
			t.Fail()
		}
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		if !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if state.GS.Gambit.Finish() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if !player.IsPlaying {
				continue
			}
			if len(player.Cards) != 3 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(state.GS.FinishRoundTime.Sub(state.GS.StartRoundTime))
		if state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if !state.GS.Gambit.Finish() {
			t.Fail()
		}
		// p1.Print()
		// p2.Print()
		// p3.Print()
	})
}

func TestLoop3(t *testing.T) {
	t.Run("when someone take check action", func(t *testing.T) {
		decisionTime := 3
		delay := 0
		minimumBet := 10
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   minimumBet}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		// dumb player
		handler.Sit("player1", 2) // dealer
		handler.Sit("player2", 5)
		handler.Sit("player3", 3) // first
		if state.GS.Gambit.Start() {
			t.Fail()
		}
		handler.StartTable()
		if !state.GS.Gambit.Start() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if len(player.Cards) != 2 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		if handler.Check("player1") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player3") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player2") {
			t.Fail()
		}
		// cannot check if already checked
		if handler.Check("player2") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player1") || !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if !player.IsPlaying {
				continue
			}
			if len(player.Cards) != 3 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(time.Second * time.Duration(delay+decisionTime))
		if handler.Check("player3") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player2") || handler.Check("player3") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player1") {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
		}
		if !state.GS.Gambit.Start() {
			t.Fail()
		}
		// _, p1 := util.Get(state.GS.Players, "player1")
		// _, p2 := util.Get(state.GS.Players, "player2")
		// _, p3 := util.Get(state.GS.Players, "player3")
		// p1.Print()
		// p2.Print()
		// p3.Print()
		// fmt.Println("now:", time.Now().Unix())
		// fmt.Println("end:", state.GS.FinishRoundTime.Unix())
	})
}

func TestLoop4(t *testing.T) {
	t.Run("when someone take bet and fold action", func(t *testing.T) {
		decisionTime := 3
		delay := 0
		minimumBet := 10
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   minimumBet}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		// dumb player
		handler.Sit("player1", 2) // dealer
		handler.Sit("player2", 4)
		handler.Sit("player3", 3) // first
		handler.Sit("player4", 5)
		handler.StartTable()
		if !state.GS.Gambit.Start() || state.GS.Pots[0] != 40 {
			t.Fail()
		}
		for _, player := range state.GS.Players {
			if player.ID == "" {
				continue
			}
			if player.Bets[0] != minimumBet {
				t.Fail()
			}
			if len(player.Cards) != 2 {
				t.Fail()
			}
			if !player.IsPlaying {
				t.Fail()
			}
			if player.Default.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player3") {
			t.Fail()
		}
		_, p1 := util.Get(state.GS.Players, "player1")
		_, p2 := util.Get(state.GS.Players, "player2")
		_, p3 := util.Get(state.GS.Players, "player3")
		_, p4 := util.Get(state.GS.Players, "player4")
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Bet("player2", 15, decisionTime) {
			t.Fail()
		}
		_, p1 = util.Get(state.GS.Players, "player1")
		_, p2 = util.Get(state.GS.Players, "player2")
		_, p3 = util.Get(state.GS.Players, "player3")
		_, p4 = util.Get(state.GS.Players, "player4")
		if p2.Bets[0] != 25 || p2.Action.Name != constant.Bet ||
			p1.Default.Name != constant.Fold ||
			p3.Default.Name != constant.Fold ||
			p4.Default.Name != constant.Fold {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Fold("player4") {
			t.Fail()
		}
		_, p4 = util.Get(state.GS.Players, "player4")
		if p4.Action.Name != constant.Fold {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay+4))
		if handler.Check("player1") || !handler.Call("player3", 3) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if handler.Check("player2") {
			t.Fail()
		}
		if !state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Bet("player3", 30, 3) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay+4))
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
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
	})
}

func TestLoop5(t *testing.T) {
	t.Run("when someone take bet and raise and fold action", func(t *testing.T) {
		decisionTime := 3
		delay := 0
		minimumBet := 10
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   minimumBet}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		// dumb player
		handler.Sit("player1", 2) // first
		handler.Sit("player2", 4)
		handler.Sit("player3", 5)
		handler.Sit("player4", 1) // dealer
		handler.StartTable()
		if !state.GS.Gambit.Start() || state.GS.Pots[0] != 40 {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Bet("player1", 20, decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Call("player2", decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay+4))
		if !handler.Bet("player4", 30, decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
			t.Fail()
		}
		if !handler.Call("player1", decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
			t.Fail()
		}
		if !handler.Bet("player2", 20, decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if handler.Call("player3", decisionTime) {
			t.Fail()
		}
		if !handler.Fold("player4") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Call("player1", decisionTime) {
			t.Fail()
		}
		if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay+3))
		if handler.Check("player1") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Bet("player2", 10, decisionTime) {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Call("player1", decisionTime) {
			t.Fail()
		}
		_, p1 := util.Get(state.GS.Players, "player1")
		_, p2 := util.Get(state.GS.Players, "player2")
		if util.SumBet(p1) != util.SumBet(p2) {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
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
	})
}

func TestLoop6(t *testing.T) {
	t.Run("when someone take bet and raise and fold action", func(t *testing.T) {
		decisionTime := 3
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
			t.Fail()
		}
		handler.Connect(id2)
		handler.Connect(id3)
		if len(state.GS.Visitors) != 3 {
			t.Fail()
		}
		handler.Sit(id1, 2) // first
		if len(state.GS.Visitors) != 2 ||
			util.CountSitting(state.GS.Players) != 1 ||
			util.CountPlaying(state.GS.Players) != 0 ||
			state.GS.Gambit.Start() {
			t.Fail()
		}
		handler.Sit(id2, 4)
		handler.Sit(id3, 5)
		if len(state.GS.Visitors) != 0 ||
			util.CountSitting(state.GS.Players) != 3 ||
			util.CountPlaying(state.GS.Players) != 0 {
			t.Fail()
		}
		handler.Connect(id4)
		if len(state.GS.Visitors) != 1 ||
			util.CountSitting(state.GS.Players) != 3 ||
			util.CountPlaying(state.GS.Players) != 0 {
			t.Fail()
		}
		handler.Sit(id4, 1) // dealer
		if len(state.GS.Visitors) != 0 ||
			util.CountSitting(state.GS.Players) != 4 ||
			util.CountPlaying(state.GS.Players) != 0 {
			t.Fail()
		}
		handler.Stand(id4)
		handler.Stand(id3)
		if len(state.GS.Visitors) != 2 ||
			util.CountSitting(state.GS.Players) != 2 ||
			util.CountPlaying(state.GS.Players) != 0 {
			t.Fail()
		}
		handler.StartTable()
		if !state.GS.Gambit.Start() {
			t.Fail()
		}
		if len(state.GS.Visitors) != 2 ||
			util.CountSitting(state.GS.Players) != 2 ||
			util.CountPlaying(state.GS.Players) != 2 {
			t.Fail()
		}
		handler.Sit(id4, 1)
		handler.Sit(id3, 5)
		if len(state.GS.Visitors) != 0 ||
			util.CountSitting(state.GS.Players) != 4 ||
			util.CountPlaying(state.GS.Players) != 2 {
			t.Fail()
		}
		if state.GS.Gambit.Start() {
			t.Fail()
		}
		if !handler.Check(id2) || !handler.Check(id1) {
			t.Fail()
		}
		if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() {
			t.Fail()
		}
		if !handler.Check(id2) || !handler.Check(id1) {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
		}
		if !state.GS.Gambit.Start() {
			t.Fail()
		}
		if !handler.Stand(id3) || len(state.GS.Visitors) != 1 ||
			util.CountSitting(state.GS.Players) != 3 ||
			util.CountPlaying(state.GS.Players) != 3 {
			t.Fail()
		}
		if !handler.Fold(id4) || !handler.Fold(id1) {
			t.Fail()
		}
		if !handler.Sit(id3, 3) {
			t.Fail()
		}
		_, player3 := util.Get(state.GS.Players, id3)
		if player3.IsPlaying || len(player3.Cards) > 0 {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
		}
		_, player2 := util.Get(state.GS.Players, id2)
		if !player2.IsWinner {
			t.Fail()
		}
		if !state.GS.Gambit.Start() ||
			!handler.Check(id1) || !handler.Check(id3) ||
			!handler.Bet(id2, 20, decisionTime) ||
			!handler.Call(id4, decisionTime) ||
			!handler.Fold(id1) || !handler.Fold(id3) {
			t.Fail()
		}
		if state.GS.Gambit.Finish() || !state.GS.Gambit.NextRound() ||
			!handler.Check(id2) || !handler.Bet(id4, 20, decisionTime) ||
			!handler.Stand(id2) {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
		}
		_, player4 := util.Get(state.GS.Players, id4)
		if !player4.IsWinner {
			t.Fail()
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
		// fmt.Println("fin:", state.GS.FinishRoundTime.Unix())
		// for _, player := range state.GS.Players {
		// 	player.Print()
		// }
	})
}
