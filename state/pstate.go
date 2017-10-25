package state

import (
	"999k_engine/model"
)

// PlayerState reperesent what each player should know
type PlayerState struct {
	Player          model.Player  `json:"player"`
	Visitors        model.Players `json:"visitors"`
	Competitors     model.Players `json:"competitors"`
	Slot            int           `json:"slot"`
	Version         int           `json:"version"`
	Pots            []int         `json:"pots"`
	CurrentTime     int64         `json:"current_time"`
	StartRoundTime  int64         `json:"start_round_time"`
	FinishRoundTime int64         `json:"finish_round_time"`
}

// Resp is server response
type Resp struct {
	Header    Header      `json:"header"`
	Payload   RespPayload `json:"payload"`
	Signature Signature   `json:"signature"`
}

// Req is server request
type Req struct {
	Header    Header
	Payload   ReqPayload
	Signature Signature
}

// Header is about token, non game logic
type Header struct {
	Token    string `json:"token"`
	DeviceID string `json:"deviceid"`
}

// RespPayload is response payload
type RespPayload struct {
	EventName string        `json:"eventname"`
	Actions   model.Actions `json:"actions"`
	GameState PlayerState   `json:"gamestate"`
}

// ReqPayload request payload
type ReqPayload struct {
	Name       string                  `json:"name"`
	Parameters model.RequestParameters `json:"parameters"`
}

// Signature is about security
type Signature struct {
}
