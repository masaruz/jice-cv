package gambit

import (
	"999k_engine/constant"
	"999k_engine/engine"
	"os"
	"strconv"
)

// Create factory of game
func Create(gambit string) engine.Gambit {
	max, err := strconv.Atoi(os.Getenv(constant.MaxPlayers))
	if err != nil {
		max = 6
	}
	dtime, err := strconv.ParseInt(os.Getenv(constant.DecisionTime), 10, 64)
	if err != nil {
		dtime = 15
	}
	mbet, err := strconv.Atoi(os.Getenv(constant.MinimumBet))
	if err != nil {
		mbet = 2
	}
	switch gambit {
	default:
		return NineK{
			MaxPlayers:   max,
			MaxAFKCount:  3,
			DecisionTime: dtime,
			MinimumBet:   mbet}
	}
}
