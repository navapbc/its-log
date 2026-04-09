package types

type SequenceEntry struct {
	Name   string         `json:"name" binding:"required"`
	Params map[string]any `json:"params"`
}
