package api

import (
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"time"
)

// Player response from api server
type Player struct {
	ID           string   `json:"userid"`
	AltID        string   `json:"alt_id"`
	DisplayName  string   `json:"display_name"`
	Clubs        []string `json:"clubs"`
	PendingClubs []string `json:"pending_clubs"`
}

// SendSticker send sticker by using gems
func SendSticker(id string) ([]byte, error) {
	// cast param to byte
	data, err := json.Marshal(Player{
		ID: id})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/sendsticker", getTableURL(id))
	return post(url, data)
}

// BuyIn when player about to sitting to table
func BuyIn(userid string, buyinamount int) ([]byte, error) {
	// Not allow empty string or buyinamount value must > 0
	if userid == "" || buyinamount <= 0 {
		return nil, fmt.Errorf("user_id or buy_in_amount is empty")
	}
	// cast param to byte
	data, err := json.Marshal(struct {
		UserID      string `json:"userid"`
		GroupID     string `json:"groupid"`
		CreateTime  int64  `json:"createtime"`
		GameIndex   int    `json:"gameindex"`
		BuyInAmount int    `json:"buyinamount"`
	}{
		UserID:      userid,
		GroupID:     state.GS.GroupID,
		CreateTime:  time.Now().Unix(),
		GameIndex:   state.GS.GameIndex,
		BuyInAmount: buyinamount,
	})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/buyin", getTableURL(state.GS.TableID))
	return post(url, data)
}

// CheckBuyIn if player is able to buyin
func CheckBuyIn(userid string, buyinamount int) ([]byte, error) {
	// Not allow empty string or buyinamount value must > 0
	if userid == "" || buyinamount <= 0 {
		return nil, fmt.Errorf("user_id or buy_in_amount is empty")
	}
	// cast param to byte
	data, err := json.Marshal(struct {
		UserID      string `json:"userid"`
		GroupID     string `json:"groupid"`
		CreateTime  int64  `json:"createtime"`
		GameIndex   int    `json:"gameindex"`
		BuyInAmount int    `json:"buyinamount"`
	}{
		UserID:      userid,
		GroupID:     state.GS.GroupID,
		CreateTime:  time.Now().Unix(),
		GameIndex:   state.GS.GameIndex,
		BuyInAmount: buyinamount,
	})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/checkbuyin", getTableURL(state.GS.TableID))
	return post(url, data)
}

// CashBack when player stand or leave the table will gain cash back from buyin
func CashBack(userid string) ([]byte, error) {
	if userid == "" {
		return nil, fmt.Errorf("user_id is empty")
	}
	_, player := util.Get(state.GS.Players, userid)
	// cast param to byte
	data, err := json.Marshal(struct {
		UserID         string  `json:"userid"`
		GroupID        string  `json:"groupid"`
		CreateTime     int64   `json:"createtime"`
		GameIndex      int     `json:"gameindex"`
		CashBackAmount float64 `json:"cashbackamount"`
	}{
		UserID:         userid,
		GroupID:        state.GS.GroupID,
		CreateTime:     time.Now().Unix(),
		GameIndex:      state.GS.GameIndex,
		CashBackAmount: player.Chips,
	})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/cashback", getTableURL(state.GS.TableID))
	return post(url, data)
}

// ExtendActionTime when player decide to extend action time
func ExtendActionTime(playerid string) ([]byte, error) {
	// TODO need to be done
	// cast param to byte
	data, err := json.Marshal(struct {
		ConsumeAmount    int    `json:"consumeamount"`
		ExtendActionTime int    `json:"extendactiontime"`
		ActionCount      int    `json:"actioncount"`
		UserID           string `json:"userid"`
		GroupID          string `json:"groupid"`
		CreateTime       int64  `json:"createtime"`
		GameIndex        int    `json:"gameindex"`
		CashBackAmount   int    `json:"cashbackamount"`
	}{
		UserID:    playerid,
		GroupID:   state.GS.GroupID,
		GameIndex: state.GS.GameIndex,
	})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/extendactiontime", getTableURL(state.GS.TableID))
	return post(url, data)
}

// ExtendTableTime extend table time
func ExtendTableTime(id string) {

}

// RemoveAuth remove player from the table
func RemoveAuth(playerid string) ([]byte, error) {
	// cast param to byte
	data, err := json.Marshal(struct {
		UserID string `json:"userid"`
	}{
		UserID: playerid,
	})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/unauth", getTableURL(state.GS.TableID))
	return post(url, data)
}
