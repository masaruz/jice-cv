package state

import (
	"999k_engine/engine"
	"999k_engine/model"
	"time"
)

// GameState to record current game
type GameState struct {
	Deck            model.Deck
	Visitors        model.Players
	Players         model.Players
	Version         int
	Turn            int
	IsTableStart    bool
	IsGameStart     bool
	StartTableTime  time.Time
	FinishTableTime time.Time
	StartRoundTime  time.Time
	FinishRoundTime time.Time
	Gambit          engine.Gambit
	Event           string
}

// gameStates to record all gamestates
var gameStates []GameState

// GS is global variable of GameState can be accessed from any where
var GS = GameState{}

// Save a gstate to gamestates
func (gstate GameState) Save() {
	// gameStates = append(gameStates, gstate)
	GS.Version++
}
