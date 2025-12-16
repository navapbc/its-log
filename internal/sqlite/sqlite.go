package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"hash/maphash"
	"log"
	"time"
	"unsafe"

	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/sqlite/models"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

type SqliteStorage struct {
	Path    string          // where the SQLite file lives
	queries *models.Queries // sqlc queries
	h       maphash.Hash    // For hashing str to int, consistently
	db      *sql.DB
}

func (s *SqliteStorage) Init() error {
	ctx := context.Background()

	t := time.Now()
	name := fmt.Sprintf("%s/%s.sqlite", s.Path, t.Format("2006-01-02"))

	db, err := sql.Open("sqlite", name)
	// FIXME: Should I create this if it doesn't exist?
	if err != nil {
		return err
	}
	s.db = db

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	s.queries = models.New(db)

	// This pulls a constant seed and dupes the maphash
	// library into using it every run as the same seed.
	// Rethink this... can we have different seeds with different runs?
	// Perhaps not in the same DB... unless... well, it would update.
	// Need the string : int mapping, and it would be OK.
	fixedSeed := viper.GetInt("app.hash_seed")
	seed := *(*maphash.Seed)(unsafe.Pointer(&fixedSeed))
	s.h.SetSeed(seed)

	return nil
}

func hashSourceAndEvent(h maphash.Hash, source string, event string) (int64, int64) {
	h.Write([]byte(source))
	source_h := h.Sum64()
	h.Reset()
	h.Write([]byte(event))
	evt_h := h.Sum64()
	h.Reset()
	return int64(source_h), int64(evt_h)
}

func (s *SqliteStorage) Event(e *itslog.Event) (int64, error) {
	source_h, evt_h := hashSourceAndEvent(s.h, e.Source, e.Event)
	// This is an unsigned to signed conversion...
	id, err := s.queries.LogEvent(context.Background(), models.LogEventParams{
		Source: source_h,
		Event:  evt_h,
	})

	if err != nil {
		panic(err)
	}

	return id, nil
}

func (s *SqliteStorage) ManyEvents(es []*itslog.Event) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	counter := int64(0)
	qtx := s.queries.WithTx(tx)
	for _, e := range es {
		if e != nil {
			// Store all of the events in a single transaction
			source_h, evt_h := hashSourceAndEvent(s.h, e.Source, e.Event)

			var err error
			_, err = qtx.LogEvent(ctx, models.LogEventParams{
				Source: source_h,
				Event:  evt_h,
			})

			if err != nil {
				log.Println("Error in storing event:" + err.Error())
				return -1, err
			}

			// Use the transaction to update the dictionary
			// in bulk as well. Individual inserts should
			// quietly ignore conflicts. This could be optimized to only update
			// when we see a new hash value.
			err = qtx.UpdateDictionary(ctx, models.UpdateDictionaryParams{
				EventSource: e.Source,
				EventName:   e.Event,
				SourceHash:  source_h,
				EventHash:   evt_h,
			})
			if err != nil {
				log.Println("Error in storing dictionary:" + err.Error())
				return -1, err
			}

			counter += 1
		}
	}

	tx.Commit()

	return counter, nil
}
