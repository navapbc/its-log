package csp

import (
	"log"
	"time"

	"github.com/jadudm/its-log/internal/itslog"
)

type EventBuffers struct {
	Events            []*itslog.Event
	eventBufferLength int
	nextEventPtr      int
	Timeout           int
}

// Do this by value, so we can pass it down a channel,
// create a new set of buffers, and not worry about
// races on the pointered structure.
func NewEventBuffers(buffer_length int) EventBuffers {
	eb := EventBuffers{
		Events: make([]*itslog.Event, buffer_length),
	}
	eb.nextEventPtr = 0
	eb.eventBufferLength = buffer_length

	return eb
}

var hasChanged bool = false

func (eb *EventBuffers) AddEvent(e *itslog.Event) bool {
	// Warning: this must be strictly sequential; this is
	// not a parallel-safe pointer update.
	eb.Events[eb.nextEventPtr] = e
	eb.nextEventPtr += 1
	// If we have a pointer >= the length, we're full
	return eb.nextEventPtr >= eb.eventBufferLength
}

func Enqueue(ch_e_in <-chan *itslog.Event, ch_flush_out chan<- EventBuffers, buffer_length int, timeout int) {
	event_buffers := NewEventBuffers(buffer_length)
	timeout_duration := time.Duration(timeout) * time.Second
	timer := time.NewTimer(timeout_duration)
	defer timer.Stop()

	for {
		select {
		case e := <-ch_e_in:
			is_full := event_buffers.AddEvent(e)
			if is_full {
				timer.Reset(timeout_duration)
				ch_flush_out <- event_buffers
				event_buffers = NewEventBuffers(buffer_length)
				hasChanged = true
			}
		case <-timer.C:
			if hasChanged {
				log.Println("flushing event buffers")
				// Send the structure out for writing
				ch_flush_out <- event_buffers
				// Allocate a new structure here in this process
				event_buffers = NewEventBuffers(buffer_length)
				hasChanged = false
			}
			timer.Reset(timeout_duration)
		}
	}
}
