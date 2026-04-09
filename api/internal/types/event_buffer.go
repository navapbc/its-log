package types

type EventBuffer struct {
	Events            []*Event
	EventBufferLength int
	NextEventPtr      int
	Timeout           int
}

// Do this by value, so we can pass it down a channel,
// create a new set of buffers, and not worry about
// races on the pointered structure.
func NewEventBuffer(bufferLength int) EventBuffer {
	eb := EventBuffer{
		Events: make([]*Event, bufferLength),
	}
	eb.NextEventPtr = 0
	eb.EventBufferLength = bufferLength

	return eb
}

func (eb *EventBuffer) AddEvent(e *Event) bool {
	// Warning: this must be strictly sequential; this is
	// not a parallel-safe pointer update.
	eb.Events[eb.NextEventPtr] = e
	eb.NextEventPtr += 1
	// If we have a pointer >= the length, we're full
	return eb.NextEventPtr >= eb.EventBufferLength
}
