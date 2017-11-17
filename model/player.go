package model

import "fmt"

// Player in the battle
type Player struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Cards         Cards    `json:"cards"`
	Chips         int      `json:"chips"`
	Type          string   `json:"type"`
	Bets          []int    `json:"bets"`
	Slot          int      `json:"slot,omitempty"`
	Default       Action   `json:"default_action"`
	Action        Action   `json:"action"`
	Actions       []Action `json:"actions"`
	IsPlaying     bool     `json:"is_playing"`
	IsEarned      bool     `json:"is_earned,omitempty"`
	IsWinner      bool     `json:"is_winner,omitempty"`
	DeadLine      int64    `json:"deadline"`
	StartLine     int64    `json:"startline"`
	WinLossAmount int      `json:"win_loss_amount,omitempty"`
	Sticker       Sticker  `json:"send_sticker,omitempty"`
}

// Print status of p only for development
func (p Player) Print() {
	fmt.Println(p.ID, p.Cards, p.Default, p.Action,
		p.StartLine, p.DeadLine, p.Chips, p.Bets,
		p.Type, p.IsWinner, p.WinLossAmount, p.Sticker)
}

// Players in the battle
type Players []Player
