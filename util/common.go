package util

import (
	"999k_engine/model"
	"999k_engine/state"
	"log"
	"time"
)

// Absolute for int64
func Absolute(num int64) int64 {
	if num < 0 {
		num = -num
	}
	return num
}

// EPSILON for testing
const EPSILON float64 = 0.00000001

// FloatEquals check equal of floats
func FloatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

// CloneState to make sure no mutable damage
func CloneState(oldState state.GameState) state.GameState {
	newState := oldState
	newState.Players = make(model.Players, len(oldState.Players))
	newState.Visitors = make(model.Players, len(oldState.Visitors))
	newState.AFKCounts = make([]int, len(oldState.AFKCounts))
	newState.DoActions = make([]bool, len(oldState.Players))
	newState.Deck.Cards = make(model.Cards, len(newState.Deck.Cards))
	newState.Pots = make([]int, len(oldState.Pots))
	newState.Scoreboard = make([]model.Scoreboard, len(oldState.Scoreboard))
	copy(newState.Players, oldState.Players)
	copy(newState.Visitors, oldState.Visitors)
	copy(newState.AFKCounts, oldState.AFKCounts)
	copy(newState.DoActions, oldState.DoActions)
	copy(newState.Deck.Cards, oldState.Deck.Cards)
	copy(newState.Pots, oldState.Pots)
	copy(newState.Scoreboard, oldState.Scoreboard)
	return newState
}

// Print log will be diabled in production
func Print(msg ...interface{}) {
	if state.Snapshot.Env != "dev" {
		log.Println(msg)
	}
}

// Log state
func Log() {
	Print(">>>>>>>>>>>> Start Broadcasting <<<<<<<<<<<<<<<")
	Print("<<< Players >>>")
	for _, player := range state.GS.Players {
		player.Print()
	}
	Print("<<< Visitors >>>")
	for _, visitor := range state.GS.Visitors {
		visitor.Print()
	}
	Print("Current gameindex:", state.GS.GameIndex)
	Print("Now:", time.Now().Unix())
	Print("Start round time:", state.GS.StartRoundTime)
	Print("Finish round time:", state.GS.FinishRoundTime)
	Print("Start table time:", state.GS.StartTableTime)
	Print("Finish table time:", state.GS.FinishTableTime)
	Print(">>>>>>>>>>>> Done Broadcasting <<<<<<<<<<<<<<<")
}
