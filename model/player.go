package model

import "fmt"

// Player in the battle
type Player struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Cards     Cards    `json:"cards"`
	Chips     int      `json:"chips"`
	Type      string   `json:"type"`
	Bets      []int    `json:"bets"`
	Slot      int      `json:"slot"`
	Default   Action   `json:"default_action"`
	Action    Action   `json:"action"`
	Actions   []Action `json:"actions"`
	IsPlaying bool     `json:"is_playing"`
	IsEarned  bool     `json:"is_earned"`
	IsWinner  bool     `json:"is_winner"`
	DeadLine  int64    `json:"deadline"`
	StartLine int64    `json:"startline"`
}

// Print status of player only for development
func (player Player) Print() {
	fmt.Println(player.ID, player.Cards, player.Default, player.Action, player.StartLine, player.DeadLine, player.Chips, player.Bets, player.Type, player.IsWinner)
}

// Players in the battle
type Players []Player

func (p Players) Len() int           { return len(p) }
func (p Players) Less(i, j int) bool { return p[i].Chips < p[j].Chips }
func (p Players) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
