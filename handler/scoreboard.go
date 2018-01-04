package handler

import (
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"math"

	"github.com/shopspring/decimal"
)

// UpdateBuyInAmount update buyin amount to scoreboard
func UpdateBuyInAmount(player *model.Player) {
	// Update scoreboard
	// If actually buyin success
	scoreboard, index := util.GetScoreboard(player.ID)
	// If not found player in scoreboard then add them
	if index == -1 {
		state.Snapshot.Scoreboard = append(state.Snapshot.Scoreboard,
			model.Scoreboard{
				UserID:      player.ID,
				DisplayName: player.Name,
				BuyInAmount: int(math.Floor(player.Chips)),
			})
	} else {
		chips := decimal.NewFromFloat(player.Chips)
		winnings := decimal.NewFromFloat(player.WinLossAmount)
		net, _ := chips.Sub(winnings).Floor().Float64()
		if netInt := int(net); netInt > scoreboard.BuyInAmount {
			scoreboard.BuyInAmount = netInt
		}
	}
}

// UpdateWinningsAmount add amount of earning chips to scoreboard
func UpdateWinningsAmount(userid string, amount float64) {
	scoreboard, index := util.GetScoreboard(userid)
	if index != -1 {
		scoreboard.WinningsAmount += amount
	}
}
