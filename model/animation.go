package model

// ClientAnimation for client to animate something
type ClientAnimation struct {
	Dealing DealingAnimation
}

// DealingAnimation for client when deal the cards
type DealingAnimation struct {
	DealingStartTime  int64 `json:"dealing_start_time"`
	DealingFinishTime int64 `json:"dealing_finish_time"`
	DealingNumber     int   `json:"dealing_number"`
}
