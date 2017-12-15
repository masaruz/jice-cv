package model

// Scoreboard history of winloss amount of this table
type Scoreboard struct {
	UserID         string  `json:"userid"`
	DisplayName    string  `json:"display_name"`
	BuyInAmount    int     `json:"buyinamount"`    // int, Buy-in amount shown on board
	WinningsAmount float64 `json:"winningsamount"` // int, Winnings amount shown on board
}
