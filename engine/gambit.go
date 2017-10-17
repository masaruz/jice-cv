package engine

// Gambit is represent of game interface
type Gambit interface {
	// Init deck and environment variables
	Init()
	// Start game
	Start()
	// NextRound game after round by round
	NextRound() bool
	// Finish game
	Finish() bool
	// End game
	End()
	// Evaluate score
	Evaluate(cards []int) (scores []int, kind string)
}
