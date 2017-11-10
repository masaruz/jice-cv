package api

import (
	"999k_engine/constant"
	"999k_engine/model"
	"encoding/json"
	"fmt"
	"log"
)

// Club response from api server
type Club struct {
	ID          string            `json:"id"`
	AltID       string            `json:"alt_id"`
	DisplayName string            `json:"display_name"`
	Members     map[string]Member `json:"members"`
}

// Member in the club
type Member struct {
	ID          string `json:"id"`
	AltID       string `json:"alt_id"`
	DisplayName string `json:"display_name"`
	Chips       int    `json:"chip"`
	JoinDate    string `json:"join_date"`
}

// ClubResp from api server
type ClubResp struct {
	Ok    bool   `json:"ok"`
	Clubs []Club `json:"clubs"`
}

// SyncPlayer to retrieve player's data from database
func SyncPlayer(id string) model.Player {
	return model.Player{ID: id, Chips: 1000}
}

// GetPlayerFromClub get player data
func GetPlayerFromClub(clubid string, memberid string) model.Player {
	body, err := get(fmt.Sprintf("%s/clubs?ids=%s", constant.Host, clubid))
	if err != nil {
		log.Fatal(err)
		return model.Player{}
	}
	resp := &ClubResp{}
	err = json.Unmarshal(body, resp)
	if err != nil || !resp.Ok {
		log.Fatal(err)
		return model.Player{}
	}
	if len(resp.Clubs) <= 0 {
		log.Fatal("clubs not found")
		return model.Player{}
	}
	member := resp.Clubs[0].Members[memberid]
	return model.Player{ID: member.ID, Chips: member.Chips}
}
