package model

// Settlement of summerized gain or loss chips
type Settlement struct {
	UserID        string  `json:"userid" validate:"required"`
	WinLossAmount int     `json:"winlossamount" validate:"required"`
	PaidRake      float64 `json:"paidrake" validate:"required"`
}
