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

func getURL(id string) string {
	return fmt.Sprintf("%s/tables/%s", Host, id)
}

// SaveToRealtime save table state to realtime
func SaveToRealtime(id string) error {
	gambit := state.GS.Gambit
	table := Table{
		ID:            id,
		GameIndex:     state.GS.GameIndex,
		PlayersAmount: util.CountSitting(state.GS.Players),
		PlayersLimit:  gambit.GetMaxPlayers(),
		BuyInMin:      gambit.GetBuyInMin(),
		BuyInMax:      gambit.GetBuyInMax(),
		DisplayName:   state.GS.TableDisplayName,
		StartTime:     state.GS.StartTableTime}
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return err
	}
	// create url
	url := fmt.Sprintf("%s/realtime", getURL(id))
	body, err := post(url, data)
	fmt.Println(string(body))
	return err
}

// DeleteFromRealtime when delete table
func DeleteFromRealtime(id string) error {
	// create url
	url := fmt.Sprintf("%s/realtime", getURL(id))
	// create request
	_, err := delete(url)
	return err
}

// GameStart set start_time only 1st game and send game index
func GameStart(id string) error {
	table := Table{}
	table.ID = id
	table.GameIndex = state.GS.GameIndex
	if table.GameIndex == 0 {
		table.StartTime = state.GS.StartRoundTime
	}
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return err
	}
	// create url
	url := fmt.Sprintf("%s/gamestart", getURL(id))
	_, err = post(url, data)
	return err
}

// Terminate caller will terminate itself
func Terminate(id string) error {
	// create url
	url := fmt.Sprintf("%s/terminate", getURL(id))
	// create request
	_, err := post(url, nil)
	return err
}

// SaveSettlements after game's end
func SaveSettlements(id string) error {
	summary := Summary{
		CreateTime: time.Now().Unix()}
	for _, player := range state.GS.Players {
		if !player.IsPlaying {
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
		return err
	}
	// create url
	url := fmt.Sprintf("%s/settlements", getURL(id))
	_, err = post(url, data)
	return err
}

// TableEnd set endtime
func TableEnd(id string) error {
	table := Table{}
	table.ID = id
	table.EndTime = time.Now().Unix()
	table.GameIndex = state.GS.GameIndex
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return err
	}
	// create url
	url := fmt.Sprintf("%s/tableend", getURL(id))
	_, err = post(url, data)
	return err
}
