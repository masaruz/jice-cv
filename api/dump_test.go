package api_test

import (
	"999k_engine/api"
	"999k_engine/gambit"
	"999k_engine/handler"
	"testing"
)

func Test01(t *testing.T) {
	decisionTime := int64(3)
	minimumBet := 10
	ninek := gambit.NineK{
		BlindsSmall:  minimumBet,
		BlindsBig:    minimumBet,
		MaxPlayers:   6,
		MaxAFKCount:  5,
		DecisionTime: decisionTime}
	handler.SetGambit(ninek)
	body, err := api.GameStart("test")
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully start game"}` {
		t.Error(data)
	}
	body, err = api.UpdateRealtimeData("test")
	if err != nil {
		t.Error()
	}
	if data := string(body); data != `{"message":"Successfully update table realtime"}` {
		t.Error(data)
	}
}
