package csp

import (
	"log"

	"github.com/jadudm/its-log/internal/itslog"
)

func FlushBuffers(ch_flush_in <-chan EventBuffers, storage itslog.ItsLog) {
	for {
		eb := <-ch_flush_in
		// Flush everything up to the pointers!
		// for ndx := 0; ndx < eb.nextEventPtr; ndx++ {
		// 	storage.Event(eb.Events[ndx])
		// }
		_, err := storage.ManyEvents(eb.Events)
		if err != nil {
			log.Printf("Failed to write event buffer; lost %d events\n", len(eb.Events))
		}
	}
}
