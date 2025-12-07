package itslog

type ItsLog interface {
	Init() error
	Event(e *Event) (int64, error)
}

type Event struct {
	Version string
	Source  string `json:"source" binding:"required"`
	Event   string `json:"event" binding:"required"`
	Value   string `json:"value" binding:"required"`
	Type    string `json:"type" binding:"required"`
}
