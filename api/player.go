package api

import (
	"999k_engine/constant"
	"999k_engine/model"
	"encoding/json"
	"fmt"
	"log"
)

// Player response from api server
type Player struct {
	ID           string   `json:"id"`
	AltID        string   `json:"alt_id"`
	DisplayName  string   `json:"display_name"`
	Clubs        []string `json:"clubs"`
	PendingClubs []string `json:"pending_clubs"`
}

// PlayerResp from api server
type PlayerResp struct {
	Ok      bool     `json:"ok"`
	Players []Player `json:"users"`
}

// GetPlayer get player data
func GetPlayer(id string) model.Player {
	body, err := get(fmt.Sprintf("%s/users?ids=%s", constant.Host, id))
	if err != nil {
		log.Fatal(err)
		return model.Player{}
	}
	resp := &PlayerResp{}
	err = json.Unmarshal(body, resp)
	if err != nil || !resp.Ok {
		log.Fatal(err)
		return model.Player{}
	}
	// get only one player
	user := resp.Players[0]
	return model.Player{ID: user.ID}
}
