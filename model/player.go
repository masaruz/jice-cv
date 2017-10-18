package model

import "fmt"
import "time"

// Player in the battle
type Player struct {
	ID        string    `json:"id"`
	Cards     Cards     `json:"cards"`
	Chips     int       `json:"chips"`
	Type      string    `json:"type"`
	Bets      []int     `json:"bets"`
	Slot      int       `json:"slot"`
	Default   Action    `json:"default"`
	Action    Action    `json:"action"`
	Actions   []Action  `json:"actions"`
	IsPlaying bool      `json:"isPlaying"`
	IsWinner  bool      `json:"isWinner"`
	DeadLine  time.Time `json:"deadline"`
	StartLine time.Time `json:"startline"`
}

// Print status of player only for development
func (player Player) Print() {
	fmt.Println(player.ID, player.IsPlaying, player.Cards, player.Default, player.Action, player.StartLine.Unix(), player.DeadLine.Unix(), player.Bets, player.Type, player.IsWinner)
}

// Players in the battle
type Players []Player
