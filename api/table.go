package api

import (
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"time"
)

// Table response from api server
type Table struct {
	ID            string `json:"tableid" validate:"required"`
	GroupID       string `json:"groupid" validate:"required"`
	GameIndex     int    `json:"gameindex" validate:"required"`
	DisplayName   string `json:"display_name,omitempty"`
	PlayersAmount int    `json:"players_amount,omitempty"`
	PlayersLimit  int    `json:"players_limit,omitempty"`
	BuyInMin      int    `json:"buyin_min,omitempty"`
	BuyInMax      int    `json:"buyin_max,omitempty"`
	StartTime     int64  `json:"start_time,omitempty"`
	EndTime       int64  `json:"end_time,omitempty"`
	Duration      int64  `json:"duration,omitempty"`
}

// Summary summery gain or loss chips for every players
type Summary struct {
	Settlements []Settlement `json:"settlements,omitempty"`
	CreateTime  int64        `json:"createtime,omitempty"`
	GameIndex   int          `json:"gameindex"`
	GroupID     string       `json:"groupid"`
}

// Settlement of summerized gain or loss chips
type Settlement struct {
	UserID        string  `json:"userid" validate:"required"`
	WinLossAmount int     `json:"winlossamount" validate:"required"`
	PaidRake      float64 `json:"paidrake" validate:"required"`
}

// Scoreboard history of winloss amount of this table
type Scoreboard struct {
	UserID         string `json:"userid"`
	DisplayName    string `json:"display_name"`
	BuyInAmount    int    `json:"buyinamount"`    // int, Buy-in amount shown on board
	WinningsAmount int    `json:"winningsamount"` // int, Winnings amount shown on board
}

// Visitor for talk to hawkeye
type Visitor struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"display_name"`
}

func getTableURL(id string) string {
	return fmt.Sprintf("%s/tables/%s", Host, id)
}

// UpdateRealtimeData save table state to realtime
func UpdateRealtimeData() ([]byte, error) {
	// Create scoreboard
	scoreboards := []Scoreboard{}
	for _, player := range state.GS.Players {
		scoreboards = append(scoreboards,
			Scoreboard{
				UserID:         player.ID,
				DisplayName:    player.Name,
				BuyInAmount:    player.Chips,
				WinningsAmount: player.WinLossAmount,
			})
	}
	// Create visitor
	visitors := []Visitor{}
	for _, visitor := range state.GS.Visitors {
		visitors = append(visitors,
			Visitor{
				UserID:      visitor.ID,
				DisplayName: visitor.Name,
			})
	}
	// gambit := state.GS.Gambit
	realtimedata := struct {
		TableID     string        `json:"tableid"`
		GroupID     string        `json:"groupid"`
		GameIndex   int           `json:"gameindex"`
		PlayerCount int           `json:"players_count"`
		Scoreboard  *[]Scoreboard `json:"scoreboard"`
		Visitors    *[]Visitor    `json:"visitors"`
	}{
		TableID:     state.GS.TableID,
		GroupID:     state.GS.GroupID,
		GameIndex:   state.GS.GameIndex,
		PlayerCount: util.CountSitting(state.GS.Players),
		Scoreboard:  &scoreboards,
		Visitors:    &visitors,
	}
	// cast param to byte
	data, err := json.Marshal(realtimedata)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/realtime", getTableURL(state.GS.TableID))
	return post(url, data)
}

// DeleteFromRealtime when delete table
func DeleteFromRealtime() ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/realtime", getTableURL(state.GS.TableID))
	// create request
	return delete(url)
}

// StartGame set start_time only 1st game and send game index
func StartGame() ([]byte, error) {
	table := Table{}
	table.GroupID = state.GS.GroupID
	table.GameIndex = state.GS.GameIndex
	if table.GameIndex == 0 {
		table.StartTime = time.Now().Unix()
	}
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/gamestart", getTableURL(state.GS.TableID))
	return post(url, data)
}

// Terminate caller will terminate itself
func Terminate() ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/terminate", getTableURL(state.GS.TableID))
	// create request
	return post(url, nil)
}

// SaveSettlements after game's end
func SaveSettlements() ([]byte, error) {
	summary := Summary{
		CreateTime: time.Now().Unix(),
		GameIndex:  state.GS.GameIndex,
		GroupID:    state.GS.GroupID,
	}
	for _, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		summary.Settlements = append(summary.Settlements, Settlement{
			UserID:        player.ID,
			WinLossAmount: player.WinLossAmount,
			PaidRake:      state.GS.Rakes[player.ID]})
	}
	// cast param to byte
	data, err := json.Marshal(summary)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/settlements", getTableURL(state.GS.TableID))
	return post(url, data)
}

// SaveSettlement support a player
func SaveSettlement(userid string) ([]byte, error) {
	_, player := util.Get(state.GS.Players, userid)
	summary := Summary{
		CreateTime: time.Now().Unix(),
		GameIndex:  state.GS.GameIndex,
		GroupID:    state.GS.GroupID,
	}
	summary.Settlements = append(summary.Settlements, Settlement{
		UserID:        player.ID,
		WinLossAmount: player.WinLossAmount,
		PaidRake:      state.GS.Rakes[player.ID]})
	// cast param to byte
	data, err := json.Marshal(summary)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/settlements", getTableURL(state.GS.TableID))
	return post(url, data)
}

// TableEnd set endtime
func TableEnd() ([]byte, error) {
	table := Table{}
	table.ID = state.GS.TableID
	table.GroupID = state.GS.GroupID
	table.EndTime = time.Now().Unix()
	table.GameIndex = state.GS.GameIndex
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/tableend", getTableURL(state.GS.TableID))
	return post(url, data)
}
