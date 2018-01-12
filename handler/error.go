package handler

import "999k_engine/model"

// ErrorCode short version
type ErrorCode = model.ErrorCode

// Error constant
const (
	ChipIsNotEnough     ErrorCode = 100
	BuyInError          ErrorCode = 101
	UpdateRealtimeError ErrorCode = 102
	CashbackError       ErrorCode = 103
	NoAvailableSeat     ErrorCode = 104
	NearOtherPlayers    ErrorCode = 105
	NoAvailableGPS      ErrorCode = 106
)
