package model

// History of hands of players
type History struct {
	Player      *PlayerHistory   `json:"player"`
	Competitors []*PlayerHistory `json:"competitors"`
}

// PlayerHistory filtered player's attributes
type PlayerHistory struct {
	ID            string  `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	WinLossAmount float64 `json:"win_loss_amount,omitempty"`
	Slot          int     `json:"slot,omitempty"`
	Cards         Cards   `json:"cards,omitempty"`
}
