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
}

// Settlement of summerized gain or loss chips
type Settlement struct {
	UserID        string  `json:"userid" validate:"required"`
	WinLossAmount int     `json:"winlossamount" validate:"required"`
	PaidRake      float64 `json:"paidrake,omitempty"`
}

func getTableURL(id string) string {
	return fmt.Sprintf("%s/tables/%s", Host, id)
}

// UpdateRealtimeData save table state to realtime
func UpdateRealtimeData() ([]byte, error) {
	gambit := state.GS.Gambit
	table := Table{
		GroupID:       state.GS.GroupID,
		GameIndex:     state.GS.GameIndex,
		PlayersAmount: util.CountSitting(state.GS.Players),
		PlayersLimit:  gambit.GetSettings().MaxPlayers,
		BuyInMin:      gambit.GetSettings().BuyInMin,
		BuyInMax:      gambit.GetSettings().BuyInMax,
		DisplayName:   state.GS.TableDisplayName,
		StartTime:     state.GS.StartTableTime}
	// cast param to byte
	data, err := json.Marshal(table)
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
	summary := Summary{CreateTime: time.Now().Unix()}
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
	summary := Summary{CreateTime: time.Now().Unix()}
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
