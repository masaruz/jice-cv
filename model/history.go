package model

// History of hands of players
type History struct {
	Player      PlayerHistory   `json:"player"`
	Competitors []PlayerHistory `json:"competitors"`
}

// PlayerHistory filtered player's attributes
type PlayerHistory struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	WinLossAmount float64 `json:"win_loss_amount"`
	Slot          int     `json:"slot"`
	CardAmount    int     `json:"card_amount"`
	Cards         Cards   `json:"cards"`
}
