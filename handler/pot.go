package handler

import (
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
)

// InitPots create and allowcate memory to handle pots
func InitPots(gs *state.GameState) {
	gs.Pots = make([]model.Pot, 5)
	for i := range gs.Pots {
		gs.Pots[i].Players = make(map[string]bool)
	}
}

// CalculatePot value and separate tier by player's bet
func CalculatePot(gs *state.GameState, id string, val int) {
	for i := range gs.Pots {
		pot := &gs.Pots[i]
		if val > pot.Ratio {
			pot.Players[id] = true
			if pot.Ratio == 0 {
				pot.Ratio = val
				break
			}
		} else if val == pot.Ratio {
			pot.Players[id] = true
			break
		} else if val < pot.Ratio {
			players := make(map[string]bool)
			for k, v := range pot.Players {
				players[k] = v
			}
			players[id] = true
			gs.Pots = append(gs.Pots[:i], append([]model.Pot{
				model.Pot{
					Ratio:   val,
					Players: players,
				},
			}, gs.Pots[i:]...)...)
			break
		}
	}
	AssignPotValue(gs)
}

// AssignPotValue base on players' bet
func AssignPotValue(gs *state.GameState) {
	prev := &model.Pot{}
	for i := range gs.Pots {
		pot := &gs.Pots[i]
		if i == 0 {
			pot.Value = pot.Ratio * Count(pot)
		} else {
			diff := pot.Ratio - prev.Ratio
			pot.Value = diff * Count(pot)
		}
		prev = pot
	}
}

// MergePots merge pot with same amount of players
func MergePots(gs *state.GameState) {
	for i := 0; i < len(gs.Pots); i++ {
		pot := &gs.Pots[i]
		for id := range pot.Players {
			_, player := util.Get(gs.Players, id)
			if player.Action.Name == constant.Fold {
				delete(pot.Players, id)
			}
		}
	}
	for i := 0; i < len(gs.Pots)-1; i++ {
		for j := i + 1; j < len(gs.Pots); j++ {
			if gs.Pots[j].Ratio == 0 {
				break
			}
			// Check players
			check, count := len(gs.Pots[i].Players), 0
			for k1 := range gs.Pots[i].Players {
				for k2 := range gs.Pots[j].Players {
					if k1 == k2 || k1 == "" {
						count++
						break
					}
				}
			}
			if check == count {
				gs.Pots[i].Value += gs.Pots[j].Value
				gs.Pots[i].Ratio = gs.Pots[j].Ratio
				gs.Pots[i].Players = gs.Pots[j].Players
				gs.Pots = append(gs.Pots[:j], gs.Pots[j+1:]...)
				break
			}
		}
	}
}

// AssignWinnerToPots who receive which pot
func AssignWinnerToPots(gs *state.GameState, id string) {
	total := util.SumPots(gs.PlayerPots)
	rakes := (gs.Gambit.GetSettings().Rake / 100) * float64(total)
	cap := gs.Gambit.GetSettings().Cap * float64(gs.Gambit.GetSettings().BlindsBig)
	if rakes > cap {
		rakes = cap
	}
	for i := range gs.Pots {
		pot := &gs.Pots[i]
		ratio := float64(pot.Value) / float64(total)
		pot.Rake = rakes * ratio
	}
	for i := len(gs.Pots) - 1; i >= 0; i-- {
		pot := &gs.Pots[i]
		if pot.Value == 0 {
			continue
		}
		for key := range pot.Players {
			if id == key {
				index, _ := util.Get(gs.Players, key)
				player := &gs.Players[index]
				related := []int{}
				for k := range pot.Players {
					slot, _ := util.Get(gs.Players, k)
					related = append(related, slot)
				}
				pot.RelatedPlayerSlots = related
				pot.Players = map[string]bool{key: true}
				pot.WinnerSlot = player.Slot
				// Receive changed
				net := 0.0
				if len(related) == 1 {
					net = float64(pot.Value)
					player.Chips += net
					state.Snapshot.Rakes[key] -= (pot.Rake / float64(len(related)))
					if state.Snapshot.Rakes[key] < 0 {
						state.Snapshot.Rakes[key] = 0
					}
					player.WinLossAmount += net
					// Is a real winner
				} else {
					net = float64(pot.Value) - pot.Rake
					player.Chips += net
					player.WinLossAmount += net
					player.IsWinner = true
				}
				AddScoreboardWinAmount(player.ID, net)
				player.IsEarned = true
				break
			}
		}
	}
}

// SumPots all values in pot
func SumPots(gs *state.GameState) int {
	sum := 0
	for _, pot := range gs.Pots {
		sum += pot.Value
	}
	return sum
}

// Count player who bet in pot for each tier
func Count(pot *model.Pot) int {
	count := 0
	for _, player := range pot.Players {
		if player {
			count++
		}
	}
	return count
}