package model

// Action what player can do
type Action struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
	Hints      Hints      `json:"Hints"`
}

// Actions list of actions that player can do
type Actions []Action

// Parameter action value
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Hint default player input value
type Hint struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// Parameters action values
type Parameters []Parameter

// Hints action value
type Hints []Hint
