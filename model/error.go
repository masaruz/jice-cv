package model

// Error handle error code and message to communicate with client
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
}

// ErrorCode handler all codes
type ErrorCode int
