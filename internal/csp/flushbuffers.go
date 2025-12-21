package csp

import (
	"log"

	"github.com/jadudm/its-log/internal/itslog"
)

// Broken out for testing
func FlushBuffersOnce(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	// Open the storage for writing before we flush
	err := storage.Init()
	if err != nil {
		panic(err)
	}

	eb := <-ch_flush_in
	_, err = storage.ManyEvents(eb.Events)
	if err != nil {
		// FIXME: really, this should percolate up to a 5xx error
		// going back to the client.
		log.Printf("Failed to write event buffer; lost %d events\n", len(eb.Events))
	}

	// Close after every flush so the DB updates/etc.
	storage.Close()
}

// For use in infinite contexts
func FlushBuffers(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {

	for {
		FlushBuffersOnce(ch_flush_in, storage)
	}
}
