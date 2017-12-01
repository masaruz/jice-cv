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

// Log state
func Log() {
	log.Println(">>>>>>>>>>>> Start Broadcasting <<<<<<<<<<<<<<<")
	log.Println("<<< Players >>>")
	for _, player := range state.GS.Players {
		player.Print()
	}
	log.Println("<<< Visitors >>>")
	for _, visitor := range state.GS.Visitors {
		visitor.Print()
	}
	log.Println("Current gameindex:", state.GS.GameIndex)
	log.Println("Now:", time.Now().Unix())
	log.Println("Start round time:", state.GS.StartRoundTime)
	log.Println("Finish round time:", state.GS.FinishRoundTime)
	log.Println("Start table time:", state.GS.StartTableTime)
	log.Println("Finish table time:", state.GS.FinishTableTime)
	log.Println(">>>>>>>>>>>> Done Broadcasting <<<<<<<<<<<<<<<")
}
