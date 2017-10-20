package model

import "fmt"

// Player in the battle
type Player struct {
	ID        string   `json:"id"`
	Cards     Cards    `json:"cards"`
	Chips     int      `json:"chips"`
	Type      string   `json:"type"`
	Bets      []int    `json:"bets"`
	Slot      int      `json:"slot"`
	Default   Action   `json:"default_action"`
	Action    Action   `json:"action"`
	Actions   []Action `json:"actions"`
	IsPlaying bool     `json:"is_playing"`
	IsWinner  bool     `json:"is_winner"`
	DeadLine  int64    `json:"deadline"`
	StartLine int64    `json:"startline"`
}

// Print status of player only for development
func (player Player) Print() {
	fmt.Println(player.ID, player.IsPlaying, player.Cards, player.Default, player.Action, player.StartLine, player.DeadLine, player.Chips, player.Bets, player.Type, player.IsWinner, player.Slot, player.Actions)
}

// Players in the battle
type Players []Player
