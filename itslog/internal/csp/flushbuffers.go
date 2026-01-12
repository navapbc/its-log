package csp

import (
	"log"

	"github.com/jadudm/its-log/internal/itslog"
)

// Broken out for testing
func FlushBuffersOnce(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	// Open the storage for writing before we flush

	eb := <-ch_flush_in
	// For testing, we get events with random dates.
	// Sort them into arrays based on date.
	// This is a waste in production.
	by_date := make(map[string][]*itslog.Event)
	for _, e := range eb.Events {
		if e != nil {
			d := e.Timestamp
			df := d.Format("2006-01-02")
			if len(by_date[df]) < 1 {
				by_date[df] = make([]*itslog.Event, 0)
			}
			by_date[df] = append(by_date[df], e)
		}
	}

	for df, es := range by_date {
		err := storage.InitByName(df)
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
