package model

// TopUp handle when player has request for topup
type TopUp struct {
	IsRequest bool    `json:"isRequest"`
	Amount    float64 `json:"amount"`
}
