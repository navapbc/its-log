package itslog

type EventType int

type ItsLog interface {
	Init() error
	Event(e *Event) (int64, error)
	ManyEvents(e []*Event) (int64, error)
}

type Event struct {
	Source string
	Event  string
}
