package api_test

import (
	"999k_engine/api"
	"999k_engine/model"
	"999k_engine/state"
	"encoding/json"
	"testing"
)

func TestLoop01(t *testing.T) {
	id1 := "us3xq4zomamja85xwx1"
	id2 := "ustel9kvy19hajahvzo3r"
	state.GS.TableID = "ta3xq4zohckj9zjaj8d"
	state.GS.GroupID = "cl3xq4zo7zojac32ay3"
	state.GS.GameIndex = 2
	state.GS.Players = model.Players{
		model.Player{
			ID:            id1,
			WinLossAmount: 100},
		model.Player{
			ID:            id2,
			WinLossAmount: -100}}
	state.GS.Rakes = map[string]float64{id1: 0.5, id2: 0.5}
	// Cannot save the settlements if players never buyin
	body, err := api.SaveSettlements()
	if err != nil {
		t.Error()
	}
	resp := &api.Response{}
	if json.Unmarshal(body, resp); resp.Error.StatusCode != 422 {
		t.Error()
	}
	body, err = api.BuyIn(id1, 200)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully buyin"}` {
		t.Error(data)
	}
	body, err = api.BuyIn(id2, 200)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully buyin"}` {
		t.Error(data)
	}
	// After they buyin success then able to save the settlement
	body, err = api.SaveSettlements()
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully settlements"}` {
		t.Error(data)
	}
	body, err = api.CashBack(id1)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully cashback"}` {
		t.Error(data)
	}
	body, err = api.CashBack(id2)
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully cashback"}` {
		t.Error(data)
	}
}
