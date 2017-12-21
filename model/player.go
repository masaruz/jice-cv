package model

import "log"

// Player in the battle
type Player struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Picture       string     `json:"picture"`
	Cards         Cards      `json:"cards"`
	CardAmount    int        `json:"card_amount"`
	Chips         float64    `json:"chips"`
	WinLossAmount float64    `json:"win_loss_amount,omitempty"`
	Type          string     `json:"type"`
	Bets          []int      `json:"bets"`
	Slot          int        `json:"slot,omitempty"`
	Default       Action     `json:"default_action"`
	Action        Action     `json:"action"`
	Actions       []Action   `json:"actions"`
	IsPlaying     bool       `json:"is_playing"`
	IsEarned      bool       `json:"is_earned,omitempty"`
	IsWinner      bool       `json:"is_winner,omitempty"`
	DeadLine      int64      `json:"deadline"`
	StartLine     int64      `json:"startline"`
	Stickers      *[]Sticker `json:"send_stickers,omitempty"`
}

// PlayerTableKey for player authentication each table
type PlayerTableKey struct {
	TableKey        string `json:"tablekey"`
	UserID          string `json:"userid"`
	ClubMemberLevel int    `json:"club_member_level"`
}

// Print status of p only for development
func (p Player) Print() {
	log.Println(p.ID, p.Name, p.Cards, p.Default, p.Action,
		p.StartLine, p.DeadLine, p.Chips, p.Bets, p.Type,
		p.IsWinner, p.WinLossAmount, p.IsPlaying, p.Stickers)
}

// Players in the battle
type Players []Player
