package csp

import (
	"log"

	"github.com/jadudm/its-log/internal/fsdb"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

// Using a nested set of maps, we end up with a structure
// that looks roughly like the following:
// app1
//
//	|
//	| - 2026-01-01 <- [ e1, e2 ]
//	| - 2026-01-02 <- [ e1, e2, e3 ]
//
// app2
//
//	\
//	 2026-01-01 <- [ e1 ]
type BufferTree map[string]map[string][]*itslog.Event

// We could be getting things from any number of apps at any given time.
// Therefore, the buffer needs to be organized for writing.
// It is a multi-level hash.
// org[appId][date] = []events
func organizeEvents(eventBuffer EventBuffers) BufferTree {
	org := make(BufferTree)

	for _, evt := range eventBuffer.Events {
		// For each event
		if evt != nil {
			// If it isn't nil, lets grab the date from the event.
			// From test endpoints, this can vary widely.
			d := evt.Timestamp
			formatted_date := d.Format("2006-01-02")

			if _, ok := org[evt.AppId]; !ok {
				// If we have not seen this app before, we have to initialize it
				org[evt.AppId] = make(map[string][]*itslog.Event)
			}

			// Now, have we seen this date before from that app?
			if _, ok := org[evt.AppId]; !ok {
				org[evt.AppId][formatted_date] = make([]*itslog.Event, 0)
			}

			// We're ready; append the event
			org[evt.AppId][formatted_date] = append(org[evt.AppId][formatted_date], evt)
		}
	}

	return org

}

func FlushBuffersOnce(s *fsdb.SqliteStorage, ch_flush_in <-chan EventBuffers) {
	eventBuffer := <-ch_flush_in
	org := organizeEvents(eventBuffer)

	for appId, dateMap := range org {
		for formatted_date, events := range dateMap {
			if s == nil {
				s = &fsdb.SqliteStorage{
					Path:     viper.GetString("storage.path"),
					Filename: formatted_date + ".sqlite",
					AppId:    appId,
					Kind:     fsdb.NamedDatabase,
				}
			}

			err := s.Init()
			if err != nil {
				panic(err)
			}
			_, err = s.ManyEvents(events)
			if err != nil {
				// FIXME: really, this should percolate up to a 5xx error
				// going back to the client. But, we don't have a Gin context, and
				// we're far away from the point where the event was logged.
				// There's no direct communication back to the client at this point, because
				// we buffered the event(s), and then flushed the buffer. This may have
				// to just be a log that we look for.
				log.Printf("Failed to write event buffer; lost %d events\n", len(events))
			}
			s.Close()
		}
	}

}

// For use in infinite contexts
func FlushBuffers(ch_flush_in <-chan EventBuffers) {
	for {
		FlushBuffersOnce(nil, ch_flush_in)
	}
}
