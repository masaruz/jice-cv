package state

import (
	"999k_engine/engine"
	"999k_engine/model"
)

// GameState to record current game
type GameState struct {
	GameIndex        int
	TableID          string
	GroupID          string
	Deck             model.Deck
	Visitors         model.Players
	Players          model.Players
	Version          int
	Turn             int
	IsTableStart     bool
	IsGameStart      bool
	IsTableExpired   bool
	StartTableTime   int64
	FinishTableTime  int64
	StartRoundTime   int64
	FinishRoundTime  int64
	MinimumBet       int
	MaximumBet       int
	Gambit           engine.Gambit
	Pots             []int
	AFKCounts        []int
	DoActions        []bool
	Event            string
	TableDisplayName string
	Rakes            map[string]float64
	PlayerTableKeys  map[string]string // Map of each player_table_key and player_id
	Scoreboard       []model.Scoreboard
}

// gameStates to record all gamestates
var gameStates []GameState

// GS is global variable of GameState can be accessed from any where
var GS = GameState{
	TableID:         "from_manager",
	GroupID:         "from_manager",
	GameIndex:       0,
	PlayerTableKeys: make(map[string]string),
}

// Snapshot is temporary gamestate used for handle state before end the script
var Snapshot = GameState{
	TableID:         "from_manager",
	GroupID:         "from_manager",
	GameIndex:       0,
	PlayerTableKeys: make(map[string]string),
}

// IncreaseVersion validate with client
func (gstate GameState) IncreaseVersion() {
	// gameStates = append(gameStates, gstate)
	GS.Version++
}
