package api

import (
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"time"
)

// Visitor in realtime database
type Visitor struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"display_name"`
}

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
	Settlements []model.Settlement `json:"settlements,omitempty"`
	CreateTime  int64              `json:"createtime,omitempty"`
	GameIndex   int                `json:"gameindex"`
	GroupID     string             `json:"groupid"`
}

func getTableURL(id string) string {
	return fmt.Sprintf("%s/tables/%s", Host, id)
}

// UpdateRealtimeData save table state to realtime
func UpdateRealtimeData() ([]byte, error) {
	// Create visitor
	visitors := []Visitor{}
	for _, visitor := range state.Snapshot.Visitors {
		if visitor.ID == "" {
			continue
		}
		visitors = append(visitors,
			Visitor{
				UserID:      visitor.ID,
				DisplayName: visitor.Name,
			})
	}
	scoreboard := []model.Scoreboard{}
	if state.Snapshot.Scoreboard != nil {
		scoreboard = state.Snapshot.Scoreboard
	}
	// gambit := state.Snapshot.Gambit
	realtimedata := struct {
		TableID     string             `json:"tableid"`
		GroupID     string             `json:"groupid"`
		GameIndex   int                `json:"gameindex"`
		PlayerCount int                `json:"players_count"`
		Scoreboard  []model.Scoreboard `json:"scoreboard"`
		Visitors    []Visitor          `json:"visitors"`
	}{
		TableID:     state.Snapshot.TableID,
		GroupID:     state.Snapshot.GroupID,
		GameIndex:   state.Snapshot.GameIndex,
		PlayerCount: util.CountSitting(state.Snapshot.Players),
		Scoreboard:  scoreboard,
		Visitors:    visitors,
	}
	data, err := json.Marshal(realtimedata)
	// cast param to byte
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/realtime", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}

// DeleteFromRealtime when delete table
func DeleteFromRealtime() ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/realtime", getTableURL(state.Snapshot.TableID))
	// create request
	return delete(url)
}

// StartGame set start_time only 1st game and send game index
func StartGame() ([]byte, error) {
	table := Table{}
	table.GroupID = state.Snapshot.GroupID
	table.GameIndex = state.Snapshot.GameIndex
	if table.GameIndex == 0 {
		table.StartTime = state.Snapshot.StartTableTime
	}
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/gamestart", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}

// Terminate caller will terminate itself
func Terminate() ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/terminate", getTableURL(state.Snapshot.TableID))
	// create request
	return post(url, nil)
}

// SaveSettlements after game's end
func SaveSettlements() ([]byte, error) {
	summary := Summary{
		CreateTime: time.Now().Unix(),
		GameIndex:  state.Snapshot.GameIndex,
		GroupID:    state.Snapshot.GroupID,
	}
	for _, player := range state.Snapshot.Players {
		if !player.IsPlaying {
			continue
		}
		summary.Settlements = append(summary.Settlements, model.Settlement{
			UserID:        player.ID,
			WinLossAmount: player.WinLossAmount,
			PaidRake:      state.Snapshot.Rakes[player.ID]})
		player.WinLossAmount = 0
	}
	// cast param to byte
	data, err := json.Marshal(summary)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/settlements", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}

// SaveSettlement support a player
func SaveSettlement(userid string) ([]byte, error) {
	_, player := util.Get(state.Snapshot.Players, userid)
	summary := Summary{
		CreateTime: time.Now().Unix(),
		GameIndex:  state.Snapshot.GameIndex,
		GroupID:    state.Snapshot.GroupID,
	}
	summary.Settlements = append(summary.Settlements, model.Settlement{
		UserID:        player.ID,
		WinLossAmount: player.WinLossAmount,
		PaidRake:      state.Snapshot.Rakes[player.ID]})
	// cast param to byte
	data, err := json.Marshal(summary)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/settlements", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}

// TableEnd set endtime
func TableEnd() ([]byte, error) {
	table := Table{EndTime: time.Now().Unix()}
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/tableend", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}
