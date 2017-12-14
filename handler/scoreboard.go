package handler

import (
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"math"
)

// UpdateScoreboard update buyin amount to scoreboard
func UpdateScoreboard(player *model.Player, action string) {
	// Update scoreboard
	// If actually buyin success
	scoreboard, index := util.GetScoreboard(player.ID)
	// If not found player in scoreboard then add them
	if index == -1 {
		state.Snapshot.Scoreboard = append(state.Snapshot.Scoreboard, model.Scoreboard{
			UserID:      player.ID,
			DisplayName: player.Name,
			BuyInAmount: int(math.Floor(player.Chips)),
		})
	} else {
		chips := int(math.Floor(player.Chips))
		switch action {
		case "add":
			scoreboard.BuyInAmount += chips
		case "remove":
			scoreboard.BuyInAmount -= chips
		}
	}
}

// AddScoreboardWinAmount add amount of earning chips to scoreboard
func AddScoreboardWinAmount(userid string, amount int) {
	scoreboard, index := util.GetScoreboard(userid)
	if index != -1 {
		scoreboard.WinningsAmount += amount
	}
}
