package itslog

import "time"

type EventType int

type ItsLog interface {
	Init() error
	// Event(e *Event) (int64, error)
	ManyEvents(e []*Event) (int64, error)
	Close()
}

type Event struct {
	Timestamp time.Time
	Cluster   string
	Source    string
	Event     string
	Value     string
}
