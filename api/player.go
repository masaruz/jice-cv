package api

import (
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"log"
)

// Player response from api server
type Player struct {
	ID           string   `json:"userid"`
	AltID        string   `json:"alt_id"`
	DisplayName  string   `json:"display_name"`
	Clubs        []string `json:"clubs"`
	PendingClubs []string `json:"pending_clubs"`
}

// PlayerResponse from api server
type PlayerResponse struct {
	Ok      bool     `json:"ok"`
	Players []Player `json:"users"`
}

// GetPlayer get player data
func GetPlayer(id string) model.Player {
	body, err := get(fmt.Sprintf("%s/users?ids=%s", Host, id))
	if err != nil {
		log.Fatal(err)
		return model.Player{}
	}
	resp := &PlayerResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil || !resp.Ok {
		log.Fatal(err)
		return model.Player{}
	}
	// get only one player
	user := resp.Players[0]
	return model.Player{ID: user.ID}
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
	url := fmt.Sprintf("%s/sendsticker", getURL(id))
	return post(url, data)
}

// BuyIn when player about to sitting to table
func BuyIn(id string) {

}

// CashBack when player stand or leave the table will gain cash back from buyin
func CashBack(id string) ([]byte, error) {
	_, player := util.Get(state.GS.Players, id)
	setttlement := Settlement{
		UserID:        player.ID,
		WinLossAmount: player.WinLossAmount,
		PaidRake:      state.GS.Rakes[player.ID]}
	// cast param to byte
	data, err := json.Marshal(setttlement)
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/cashback", getURL(id))
	return post(url, data)
}

// ExtendActionTime when player decide to extend action time
func ExtendActionTime(id string) ([]byte, error) {
	// cast param to byte
	data, err := json.Marshal(Player{
		ID: id})
	if err != nil {
		return nil, err
	}
	// create url
	url := fmt.Sprintf("%s/extendactiontime", getURL(id))
	return post(url, data)
}

// ExtendTableTime extend table time
func ExtendTableTime(id string) {

}

// RemoveAuth remove player from the table
func RemoveAuth(id string) {

}
