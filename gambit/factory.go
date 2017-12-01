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
	bsm, err := strconv.Atoi(os.Getenv(constant.BlindsSmall))
	if err != nil {
		bsm = 2
	}
	bbg, err := strconv.Atoi(os.Getenv(constant.BlindsBig))
	if err != nil {
		bbg = 2
	}
	rake, err := strconv.ParseFloat(os.Getenv(constant.Rake), 64)
	if err != nil {
		rake = 5.0 // percent
	}
	cap, err := strconv.ParseFloat(os.Getenv(constant.Cap), 64)
	if err != nil {
		cap = 0.5
	}
	// Minimum buy-in
	minbi, err := strconv.Atoi(os.Getenv(constant.MinimumBuyIn))
	if err != nil {
		minbi = 200
	}
	// Maximum buy-in
	maxbi, err := strconv.Atoi(os.Getenv(constant.MaximumBuyIn))
	if err != nil {
		maxbi = 1000
	}
	switch gambit {
	default:
		return NineK{
			MaxPlayers:   max,
			MaxAFKCount:  3,
			DecisionTime: dtime,
			BlindsSmall:  bsm,
			BlindsBig:    bbg,
			Rake:         rake,
			Cap:          cap,
			BuyInMax:     maxbi,
			BuyInMin:     minbi}
	}
}
