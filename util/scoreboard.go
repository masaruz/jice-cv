package util

import (
	"999k_engine/model"
	"999k_engine/state"
)

// GetScoreboard get score board from userid
func GetScoreboard(userid string) (*model.Scoreboard, int) {
	for index, sb := range state.GS.Scoreboard {
		if sb.UserID == userid {
			return &state.GS.Scoreboard[index], index
		}
	}
	return nil, -1
}

// AddScoreboardWinAmount add amount of earning chips to scoreboard
func AddScoreboardWinAmount(userid string, amount int) {
	scoreboard, index := GetScoreboard(userid)
	if index != -1 {
		scoreboard.WinningsAmount += amount
	}
}
