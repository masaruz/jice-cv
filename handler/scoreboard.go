package handler

import (
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
)

// UpdateScoreboard update buyin amount to scoreboard
func UpdateScoreboard(player *model.Player, action string) {
	// Update scoreboard
	// If actually buyin success
	scoreboard, sbindex := util.GetScoreboard(player.ID)
	// If not found player in scoreboard then add them
	if sbindex == -1 {
		state.Snapshot.Scoreboard = append(state.Snapshot.Scoreboard, model.Scoreboard{
			UserID:      player.ID,
			DisplayName: player.Name,
			BuyInAmount: player.Chips,
		})
	} else {
		switch action {
		case "add":
			scoreboard.BuyInAmount += player.Chips
		case "remove":
			scoreboard.BuyInAmount -= player.Chips
		}
	}
}
