package game_test

import (
	"999k_websocket/constant"
	"999k_websocket/game"
	"999k_websocket/handler"
	"999k_websocket/state"
	"999k_websocket/util"
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
		if p4.DeadLine.Sub(state.GS.StartGameTime).Seconds() != float64(4*decisionTime) ||
			p3.DeadLine.Sub(state.GS.StartGameTime).Seconds() != float64(2*decisionTime) ||
			p1.DeadLine.Sub(state.GS.StartGameTime).Seconds() != float64(1*decisionTime) ||
			p2.DeadLine.Sub(state.GS.StartGameTime).Seconds() != float64(3*decisionTime) {
			t.Fail()
		}
		// nothing happend in 2 seconds and assume players act default action
		time.Sleep(time.Second * time.Duration(4))
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
		time.Sleep(time.Second * time.Duration(4))
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
		time.Sleep(time.Second * time.Duration(2))
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
		time.Sleep(time.Second * time.Duration(2))
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
		time.Sleep(time.Second * time.Duration(3))
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
		time.Sleep(time.Second * time.Duration(3))
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
	t.Run("when someone take action", func(t *testing.T) {
		decisionTime := 1
		minimumBet := 10
		ninek := game.NineK{
			MaxPlayers:   6,
			DecisionTime: decisionTime,
			MinimumBet:   minimumBet}
		handler.SetGambit(ninek)
		state.GS.Gambit.Init() // create seats
		// dumb player
		handler.Sit("player1", 2)
		handler.Sit("player2", 5)
		handler.Sit("player3", 3)
		handler.StartTable()
		state.GS.Gambit.Start()
	})
}
