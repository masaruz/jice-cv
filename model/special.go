package model

// Sticker that can be send to another
type Sticker struct {
	ID         string `json:"sticker_id,omitempty"`
	StartTime  int64  `json:"start_time,omitempty"`
	FinishTime int64  `json:"finish_time,omitempty"`
	ToTarget   int    `json:"to_target,omitempty"`
}
