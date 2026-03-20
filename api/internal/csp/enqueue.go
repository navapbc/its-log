// Package csp provides CSP-style processes for handling events
//
// CSP, or Communicating Sequential Processes, developed by Tony Hoare,
// is a formal algebra for describing and reasoning about parallel processes.
//
// https://en.wikipedia.org/wiki/Communicating_sequential_processes
//
// It was used to specify and verify the Transputer, formed the basis for
// the programming language occam, and is lies at the heart of Golang's
// channel and gofunc abstractions.
//
// Enqueue provides a process to consume events from the API, serializing them
// so we can write them in batches to the underlying SQLite database
//
// FlushBuffers handles buffered events that have been enqeueued, and writes
// them out to the underlying database.
//
// This pair of processes eliminate concerns about writing in parallel to the
// filesystem-based SQLite database underneath its-log.
package csp

import (
	"log"
	"time"

	"github.com/jadudm/its-log/internal/base"
)

type EventBuffers struct {
	Events            []*base.Event
	eventBufferLength int
	nextEventPtr      int
	Timeout           int
}

// Do this by value, so we can pass it down a channel,
// create a new set of buffers, and not worry about
// races on the pointered structure.
func NewEventBuffers(buffer_length int) EventBuffers {
	eb := EventBuffers{
		Events: make([]*base.Event, buffer_length),
	}
	eb.nextEventPtr = 0
	eb.eventBufferLength = buffer_length

	return eb
}

func (eb *EventBuffers) AddEvent(e *base.Event) bool {
	// Warning: this must be strictly sequential; this is
	// not a parallel-safe pointer update.
	eb.Events[eb.nextEventPtr] = e
	eb.nextEventPtr += 1
	// If we have a pointer >= the length, we're full
	return eb.nextEventPtr >= eb.eventBufferLength
}

func Enqueue(ch_evt_in <-chan *base.Event, ch_flush_out chan<- EventBuffers, buffer_length int, timeout int) {
	event_buffers := NewEventBuffers(buffer_length)
	timeout_duration := time.Duration(timeout) * time.Second
	timer := time.NewTimer(timeout_duration)
	defer timer.Stop()

	for {
		select {
		case e := <-ch_evt_in:
			is_full := event_buffers.AddEvent(e)
			timer.Reset(timeout_duration)
			if is_full {
				log.Println("flushing full buffers")
				ch_flush_out <- event_buffers
				event_buffers = NewEventBuffers(buffer_length)
			}
		case <-timer.C:
			// This will flush once at startup, because the timer fires.
			// This has a side-effect of creating the DB.
			log.Println("flushing stale buffers")
			// Send the structure out for writing
			ch_flush_out <- event_buffers
			// Allocate a new structure here in this process
			event_buffers = NewEventBuffers(buffer_length)
			// Do not reset the timer here. Only reset if
			// new events come through, and they might need to
			// be flushed before the buffer is full.
		}
	}
}
