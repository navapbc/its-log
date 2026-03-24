package types

type EventBuffers struct {
	Events            []*Event
	EventBufferLength int
	NextEventPtr      int
	Timeout           int
}

func (eb *EventBuffers) AddEvent(e *Event) bool {
	// Warning: this must be strictly sequential; this is
	// not a parallel-safe pointer update.
	eb.Events[eb.NextEventPtr] = e
	eb.NextEventPtr += 1
	// If we have a pointer >= the length, we're full
	return eb.NextEventPtr >= eb.EventBufferLength
}
