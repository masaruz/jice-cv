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
			if player.Action.Name != constant.Check {
				t.Fail()
			}
		}
		// test timeline
		_, p1 := util.Get(state.GS.Players, "player1")
		_, p2 := util.Get(state.GS.Players, "player2")
		_, p3 := util.Get(state.GS.Players, "player3")
		_, p4 := util.Get(state.GS.Players, "player4")
		newDecisionTime := decisionTime + 1
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
			if player.Action.Name != constant.Check {
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
		handler.StartTable()
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
			if player.Action.Name != constant.Check {
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
			if player.Action.Name != constant.Check {
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
			if player.Action.Name != constant.Check {
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
			if player.Action.Name != constant.Check {
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
		delay := 1
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
		handler.StartTable()
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
			if player.Action.Name != constant.Check {
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
		if handler.Check("player2") || handler.Check("player1") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player1") || state.GS.Gambit.NextRound() {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !state.GS.Gambit.NextRound() {
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
			if player.Action.Name != constant.Check {
				t.Fail()
			}
		}
		time.Sleep(time.Second * time.Duration(delay+decisionTime))
		if handler.Check("player3") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player2") || handler.Check("player3") || handler.Check("player1") {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if !handler.Check("player1") {
			t.Fail()
		}
		if state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
			t.Fail()
		}
		time.Sleep(time.Second * time.Duration(delay))
		if state.GS.Gambit.NextRound() || !state.GS.Gambit.Finish() {
			t.Fail()
		}
		// _, p1 := util.Get(state.GS.Players, "player1")
		// _, p2 := util.Get(state.GS.Players, "player2")
		// _, p3 := util.Get(state.GS.Players, "player3")
		// p1.Print()
		// p2.Print()
		// p3.Print()
	})
}
