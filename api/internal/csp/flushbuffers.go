package csp

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

// Using a nested set of maps, we end up with a structure
// that looks roughly like the following:
//
// app1
//
//	|
//	| - 2026-01-01 <- [ e1, e3 ]
//	| - 2026-01-02 <- [ e2, e4, e5 ]
//
// app2
//
//	\
//	 2026-01-01 <- [ e6 ]
type BufferTree map[string]map[string][]*types.Event

// We could be getting things from any number of apps at any given time.
// Therefore, the buffer needs to be organized for writing.
// It is a multi-level hash.
// org[appId][date] = []events
func organizeEvents(eventBuffer types.EventBuffer) BufferTree {
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
				org[evt.AppId] = make(map[string][]*types.Event)
			}

			// Now, have we seen this date before from that app?
			if _, ok := org[evt.AppId]; !ok {
				org[evt.AppId][formatted_date] = make([]*types.Event, 0)
			}

			// We're ready; append the event
			org[evt.AppId][formatted_date] = append(org[evt.AppId][formatted_date], evt)
		}
	}

	return org

}

var isEtlLoaded = make(map[string]bool)

func FlushBuffersOnce(ch_flush_in <-chan types.EventBuffer) {
	eventBuffer := <-ch_flush_in
	org := organizeEvents(eventBuffer)

	for appId, dateMap := range org {
		for formattedDate, events := range dateMap {
			s := types.NewStorage(appId)
			err := s.SetDate(formattedDate)
			if err != nil {
				panic("failed to parse date in FlushBuffersOnce")
			}

			err = s.Init()
			if err != nil {
				log.Println("storage init error: " + err.Error())
				panic(err)
			}

			_, err = ManyEvents(s, events)

			if err != nil {
				// FIXME: really, this should percolate up to a 5xx error
				// going back to the client. But, there's no direct communication back to
				// the client at this point, because we buffered the event(s), and then
				// flushed the buffer. This may have to just be a log that we look for.
				log.Printf("Failed to write event buffer; lost %d events\n", len(events))
				log.Println("err: " + err.Error())
			}

			// This `if` saves us from attempting to reload multiple times for a single
			// running session of its-log, but it also means there is no hot-reloading
			// of the defaults. (Or, there could be, but it would have to be an API call.)
			// I think this keeps us from hitting the DB too often, so I'm going to leave it.
			if _, ok := isEtlLoaded[formattedDate]; !ok {
				err := base.LoadDefaultEtlFiles(s)
				if err != nil {
					log.Println("could not load default ETL files for " + formattedDate + ": " + err.Error())
				} else {
					isEtlLoaded[formattedDate] = true
				}
			}
			s.Close()
		}
	}

}

// For use in infinite contexts
func FlushBuffers(ch_flush_in <-chan types.EventBuffer) {
	for {
		FlushBuffersOnce(ch_flush_in)
	}
}

func ManyEvents(s *types.Storage, evt_buff []*types.Event) (int64, error) {
	counter := int64(0)

	s.Lock()
	// Unlock regardless of whether we complete
	// or leave part way through the loop.
	defer s.Unlock()

	for _, e := range evt_buff {
		if e != nil {

			valid_cluster := false
			if len(e.Cluster) > 0 {
				valid_cluster = true
			}

			valid_value := false
			if len(e.Value) > 0 {
				valid_value = true
			}

			_, err := s.Queries.LogEvent(context.Background(), models.LogEventParams{
				Timestamp: e.Timestamp.Format(time.RFC3339),
				KeyID:     e.KeyId,
				Cluster:   sql.NullString{String: e.Cluster, Valid: valid_cluster},
				Tags:      e.TagString,
				Value:     sql.NullString{String: e.Value, Valid: valid_value},
			})

			if err != nil {
				log.Println("Error in storing event: " + err.Error())
				return -1, err
			}
			counter += 1
		}
	}

	return counter, nil
}
