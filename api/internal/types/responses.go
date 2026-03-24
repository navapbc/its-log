package types

type Success struct {
	Status string `json:"status"`
}
type Error struct {
	Status    string         `json:"status"`
	Data      map[string]any `json:"data"`
	ErrorType string         `json:"error_type"`
	Error     string         `json:"error"`
}
