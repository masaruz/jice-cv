package api

// constant use for api package
const (
	Host   = "https://xysfh0h6qc.execute-api.ap-southeast-1.amazonaws.com/dev"
	APIKey = "hozbzFOs516fH6Z5kEgwq21nHJhBhSHW6qvbvkmW"
)

// Response from api server
type Response struct {
	Error Error `json:"err"`
}

// Error when receive something wrong from server
type Error struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}
