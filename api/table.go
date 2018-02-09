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

// History post body
type History struct {
	UserID      string                `json:"userid"`
	TableID     string                `json:"tableid"`
	Name        string                `json:"name"`
	CreateTime  int64                 `json:"createtime"`
	Player      model.PlayerHistory   `json:"player"`
	Competitors []model.PlayerHistory `json:"competitors"`
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

// StartTable set start_time only 1st game and send game index
func StartTable() ([]byte, error) {
	table := Table{}
	table.StartTime = state.Snapshot.StartTableTime
	// cast param to byte
	data, err := json.Marshal(table)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/tablestart", getTableURL(state.Snapshot.TableID))
	return post(url, data)
}

// StartGame set start_time only 1st game and send game index
func StartGame() ([]byte, error) {
	table := Table{}
	table.GroupID = state.Snapshot.GroupID
	table.GameIndex = state.Snapshot.GameIndex
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

// SaveHistories save latest histories when game is end
func SaveHistories() ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/history", getTableURL(state.Snapshot.TableID))
	// cast param to byte
	histories := []History{}
	for _, history := range state.Snapshot.History {
		histories = append(histories, History{
			UserID:      history.Player.ID,
			TableID:     state.Snapshot.TableID,
			Name:        history.Player.Name,
			CreateTime:  time.Now().Unix(),
			Player:      history.Player,
			Competitors: history.Competitors,
		})
	}
	data, err := json.Marshal(struct {
		Histories []History `json:"histories"`
	}{Histories: histories})
	if err != nil {
		return nil, err
	}
	// create request
	return post(url, data)
}

// GetHistories get list of hand history
func GetHistories(userid string) ([]byte, error) {
	// create url
	url := fmt.Sprintf("%s/history?userid=%s", getTableURL(state.Snapshot.TableID), userid)
	return get(url)
}

// SaveSettlements after game's end
func SaveSettlements() ([]byte, error) {
	summary := Summary{
		CreateTime: time.Now().Unix(),
		GameIndex:  state.Snapshot.GameIndex,
		GroupID:    state.Snapshot.GroupID,
	}
	for index := range state.Snapshot.Players {
		player := &state.Snapshot.Players[index]
		if !player.IsPlaying {
			continue
		}
		summary.Settlements = append(summary.Settlements, model.Settlement{
			UserID:        player.ID,
			WinLossAmount: player.WinLossAmount,
			PaidRake:      state.Snapshot.Rakes[player.ID]})
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
