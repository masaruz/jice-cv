package engine

import (
	"999k_engine/model"
)

// Gambit is represent of game interface
type Gambit interface {
	// Init deck and environment variables
	Init()
	// Start game
	Start() bool
	// Check is doing nothing
	Check(id string) bool
	// Bet is raise chips
	Bet(id string, chips int) bool
	// Raise is raise chips
	Raise(id string, chips int) bool
	// Call is raise chips to equal highest bet chips
	Call(id string) bool
	// Fold is out of the game
	Fold(id string) bool
	// AllIn invest everything
	AllIn(id string) bool
	// NextRound game after round by round
	NextRound() bool
	// Finish game
	Finish() bool
	// End game
	End()
	// Evaluate score
	Evaluate(cards []int) (scores []int, kind string)
	// Reducer event and return action
	Reducer(event string, id string) model.Actions

	GetMaxPlayers() int
	GetBlindsSmall() int
	GetBlindsBig() int
	GetBuyInMin() int
	GetBuyInMax() int
	GetRake() float64
	GetCap() float64
	GetDecisionTime() int64
}
