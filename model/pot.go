package model

import "log"

// Pot handle situation of distributed pot
type Pot struct {
	Ratio      int
	Players    map[string]bool
	Value      int `json:"value"`
	WinnerSlot int `json:"winner_slot"`
}

// Pots are array of Pot
type Pots []Pot

// Print pots details
func (pots Pots) Print() {
	for _, pot := range pots {
		log.Printf("value=%d ratio=%d players=%v", pot.Value, pot.Ratio, pot.Players)
	}
}
