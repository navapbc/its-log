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

func Enqueue(ch_evt_in <-chan *types.Event, ch_flush_out chan<- types.EventBuffer, bufferLength int, timeout int) {
	event_buffers := types.NewEventBuffer(bufferLength)
	timeoutDuration := time.Duration(timeout) * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	for {
		select {
		case e := <-ch_evt_in:
			is_full := event_buffers.AddEvent(e)
			// We reset the timer here because we're "live" and getting
			// events in less than timeoutDuration seconds. So, let's not flush
			// yet. Flushing the buffer too frequently leads to bad performance
			// when under heavy logging load.
			timer.Reset(timeoutDuration)
			// If the buffer is full, we should flush and create a new buffer.
			// This lets us keep grabbing events from the API handler while we are
			// writing this buffer to the DB.
			if is_full {
				log.Println("flushing full buffers")
				ch_flush_out <- event_buffers
				event_buffers = types.NewEventBuffer(bufferLength)
			}
		case <-timer.C:
			// This will flush once at startup, because the timer fires.
			// This has a side-effect of creating the DB.
			log.Printf("Enqueue: flushing stale buffers\n")
			// Send the structure out for writing
			ch_flush_out <- event_buffers
			// Allocate a new structure here in this process
			event_buffers = types.NewEventBuffer(bufferLength)
			// Do not reset the timer here. Only reset if
			// new events come through, and they might need to
			// be flushed before the buffer is full.
		}
	}
}
