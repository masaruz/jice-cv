package util

import (
	"999k_engine/model"
	"999k_engine/state"
)

// GetScoreboard get score board from userid
func GetScoreboard(userid string) (*model.Scoreboard, int) {
	for index, sb := range state.Snapshot.Scoreboard {
		if sb.UserID == userid {
			return &state.Snapshot.Scoreboard[index], index
		}
	}
	return nil, -1
}
