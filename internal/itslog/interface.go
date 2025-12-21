package itslog

import "time"

type EventType int

type ItsLog interface {
	Init(t time.Time) error
	Event(e *Event) (int64, error)
	ManyEvents(e []*Event) (int64, error)
	Summarize()
	Close()
}

type Event struct {
	Source string
	Event  string
}
