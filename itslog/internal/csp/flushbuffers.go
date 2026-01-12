package csp

import (
	"log"
	"time"

	"github.com/jadudm/its-log/internal/itslog"
)

// Broken out for testing
func FlushBuffersOnce(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	// Open the storage for writing before we flush

	eb := <-ch_flush_in
	// For testing, we get events with random dates.
	// Sort them into arrays based on date.
	// This is a waste in production.
	by_date := make(map[time.Time][]*itslog.Event)
	for _, e := range eb.Events {
		if e != nil {
			d := e.Timestamp
			if len(by_date[d]) < 1 {
				by_date[d] = make([]*itslog.Event, 0)
			}
			by_date[d] = append(by_date[d], e)

		}
	}

	for d, es := range by_date {
		err := storage.InitByDate(d)
		if err != nil {
			panic(err)
		}
		_, err = storage.ManyEvents(es)
		if err != nil {
			// FIXME: really, this should percolate up to a 5xx error
			// going back to the client.
			log.Printf("Failed to write event buffer; lost %d events\n", len(es))
		}
		storage.Close()
	}
}

// For use in infinite contexts
func FlushBuffers(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {

	for {
		FlushBuffersOnce(ch_flush_in, storage)
	}
}
