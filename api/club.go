package api

// Club response from api server
type Club struct {
	ID          string            `json:"id"`
	AltID       string            `json:"alt_id"`
	DisplayName string            `json:"display_name"`
	Members     map[string]Member `json:"members"`
}

// Member in the club
type Member struct {
	ID          string `json:"id"`
	AltID       string `json:"alt_id"`
	DisplayName string `json:"display_name"`
	Chips       int    `json:"chip"`
	JoinDate    string `json:"join_date"`
}

// ClubResponse from api server
type ClubResponse struct {
	Ok    bool   `json:"ok"`
	Clubs []Club `json:"clubs"`
}
