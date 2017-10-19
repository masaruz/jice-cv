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
	Header    Header    `json:"header"`
	Payload   Payload   `json:"payload"`
	Signature Signature `json:"signature"`
}

// Header is about token, non game logic
type Header struct {
	Token string `json:"token"`
}

// Payload is gameplay data
type Payload struct {
	EventName string        `json:"eventname"`
	Actions   model.Actions `json:"actions"`
	GameState PlayerState   `json:"gamestate"`
}

// Signature is about security
type Signature struct {
}
