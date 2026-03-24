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

	"github.com/navapbc/its-log/internal/types"
)

// Do this by value, so we can pass it down a channel,
// create a new set of buffers, and not worry about
// races on the pointered structure.
func NewEventBuffers(buffer_length int) types.EventBuffers {
	eb := types.EventBuffers{
		Events: make([]*types.Event, buffer_length),
	}
	eb.NextEventPtr = 0
	eb.EventBufferLength = buffer_length

	return eb
}

func Enqueue(ch_evt_in <-chan *types.Event, ch_flush_out chan<- types.EventBuffers, bufferLength int, timeout int) {
	event_buffers := NewEventBuffers(bufferLength)
	timeoutDuration := time.Duration(timeout) * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	for {
		select {
		case e := <-ch_evt_in:
			is_full := event_buffers.AddEvent(e)
			timer.Reset(timeoutDuration)
			if is_full {
				log.Println("flushing full buffers")
				ch_flush_out <- event_buffers
				event_buffers = NewEventBuffers(bufferLength)
			}
		case <-timer.C:
			// This will flush once at startup, because the timer fires.
			// This has a side-effect of creating the DB.
			log.Println("flushing stale buffers")
			// Send the structure out for writing
			ch_flush_out <- event_buffers
			// Allocate a new structure here in this process
			event_buffers = NewEventBuffers(bufferLength)
			// Do not reset the timer here. Only reset if
			// new events come through, and they might need to
			// be flushed before the buffer is full.
		}
	}
}
