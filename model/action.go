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

// RequestParameter is from client
type RequestParameter struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ValueInteger int    `json:"value_integer"`
	ValueString  string `json:"value_string"`
}

// Parameters action values
type Parameters []Parameter

// RequestParameters list of RequestParameter
type RequestParameters []RequestParameter

// Hints action value
type Hints []Hint
