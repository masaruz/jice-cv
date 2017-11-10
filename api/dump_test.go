package api_test

import (
	"999k_engine/api"
	"testing"
)

func Test01(t *testing.T) {
	player := api.GetPlayer("us3xq4zonwqj9is76ch")
	if player.ID == "" {
		t.Error()
	}
}

func Test02(t *testing.T) {
	player := api.GetPlayerFromClub("cl3xq4zop81j9jlvo9g", "us3xq4zop81j9jltiwj")
	if player.ID == "" {
		t.Error()
	}
}
