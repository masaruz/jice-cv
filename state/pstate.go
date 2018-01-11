package state

import (
	"999k_engine/model"
)

// PlayerState reperesent what each player should know
type PlayerState struct {
	Player       model.Player  `json:"player"`
	Visitors     model.Players `json:"visitors"`
	Competitors  model.Players `json:"competitors"`
	Slot         int           `json:"slot"`
	Version      int           `json:"version"`
	Pots         []int         `json:"pots"`
	SummaryPots  model.Pots    `json:"summary_pots"`
	HighestBet   int           `json:"highest_bet"`
	IsGameStart  bool          `json:"is_game_start"`
	IsTableStart bool          `json:"is_table_start"`
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
	Token           string `json:"token"`
	DeviceID        string `json:"deviceid"`
	DisplayName     string `json:"display_name"`
	AvatarSource    string `json:"avatar_source,omitempty"`
	AvatarBuiltinID string `json:"avatar_builtin_id,omitempty"`
	AvatarCustomID  string `json:"avatar_custom_id,omitempty"`
	FacebookID      string `json:"facebook_id,omitempty"`
}

// RespPayload is response payload
type RespPayload struct {
	EventName       string             `json:"eventname"`
	Actions         model.Actions      `json:"actions"`
	GameState       PlayerState        `json:"gamestate"`
	Scoreboard      []model.Scoreboard `json:"scoreboard"`
	GameIndex       int                `json:"gameindex"`
	CurrentTime     int64              `json:"current_time"`
	StartRoundTime  int64              `json:"start_round_time"`
	FinishRoundTime int64              `json:"finish_round_time"`
	IsTableExpired  bool               `json:"is_table_expired"`
	FinishGameDelay int64              `json:"finish_game_delay"`
	Error           *model.Error       `json:"error,omitempty"`
}

// ReqPayload request payload
type ReqPayload struct {
	Name       string                  `json:"name"`
	Parameters model.RequestParameters `json:"parameters"`
}

// Signature is about security
type Signature struct {
}
