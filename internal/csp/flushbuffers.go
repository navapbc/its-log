package csp

import (
	"log"

	"github.com/jadudm/its-log/internal/itslog"
)

// Broken out for testing
func FlushBuffersOnce(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	eb := <-ch_flush_in
	_, err := storage.ManyEvents(eb.Events)
	if err != nil {
		// FIXME: really, this should percolate up to a 5xx error
		// going back to the client.
		log.Printf("Failed to write event buffer; lost %d events\n", len(eb.Events))
	}
}

// For use in infinite contexts
func FlushBuffers(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	for {
		FlushBuffersOnce(ch_flush_in, storage)
	}
}
