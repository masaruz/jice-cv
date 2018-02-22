package util

import (
	"999k_engine/model"
	"999k_engine/state"
	"log"
	"math"
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
	newState.PlayerPots = make([]int, len(oldState.PlayerPots))
	newState.Scoreboard = make([]model.Scoreboard, len(oldState.Scoreboard))
	newState.Pots = make(model.Pots, len(oldState.Pots))
	copy(newState.Players, oldState.Players)
	copy(newState.Visitors, oldState.Visitors)
	copy(newState.AFKCounts, oldState.AFKCounts)
	copy(newState.DoActions, oldState.DoActions)
	copy(newState.Deck.Cards, oldState.Deck.Cards)
	copy(newState.PlayerPots, oldState.PlayerPots)
	copy(newState.Scoreboard, oldState.Scoreboard)
	copy(newState.Pots, oldState.Pots)
	for i := range oldState.Pots {
		newState.Pots[i].Players = make(map[string]bool)
		for k, v := range oldState.Pots[i].Players {
			newState.Pots[i].Players[k] = v
		}
	}
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
	Print("Duration:", state.GS.Duration)
	Print("Finish table time:", state.GS.FinishTableTime)
	Print(">>>>>>>>>>>> Done Broadcasting <<<<<<<<<<<<<<<")
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance calculate distance between to location
func Distance(player, other model.Player) float64 {
	if player.Lon == 0 && player.Lat == 0 {
		return 0
	}
	var la1, lo1, la2, lo2, r float64
	la1 = player.Lat * math.Pi / 180
	lo1 = player.Lon * math.Pi / 180
	la2 = other.Lat * math.Pi / 180
	lo2 = other.Lon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
