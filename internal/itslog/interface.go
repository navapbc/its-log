package itslog

type ItsLog interface {
	Init() error
	Event(e *Event) (int64, error)
}

type Event struct {
	Source string `json:"source" binding:"required"`
	Event  string `json:"event" binding:"required"`
}
