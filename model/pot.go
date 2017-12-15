package model

import "log"

// Pot handle situation of distributed pot
type Pot struct {
	Ratio              int
	Rake               float64
	Players            map[string]bool
	Value              int   `json:"value"`
	WinnerSlot         int   `json:"winner_slot"`
	RelatedPlayerSlots []int `json:"related_player_slots"`
}

// Pots are array of Pot
type Pots []Pot

// Print pots details
func (pots Pots) Print() {
	for _, pot := range pots {
		log.Printf("value=%d ratio=%d rake=%v players=%v related=%v", pot.Value, pot.Ratio, pot.Rake, pot.Players, pot.RelatedPlayerSlots)
	}
}
