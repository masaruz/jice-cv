package util

import (
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

// Print log will be diabled in production
func Print(msg ...interface{}) {
	if state.GS.Env != "dev" {
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
