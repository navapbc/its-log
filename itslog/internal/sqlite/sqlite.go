package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"hash/maphash"
	"log"
	"strings"
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
	Path     string          // where the SQLite file lives
	Filename string          // an optional filename override
	queries  *models.Queries // sqlc queries
	h        maphash.Hash    // For hashing str to int, consistently
	db       *sql.DB
}

func _init(s *SqliteStorage, name string) error {
	ctx := context.Background()

	db, err := sql.Open("sqlite", name)
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
	// This seed matters for how strings are mapped to integers
	// If it changes from one run to the next, then the mapping
	// will change. That is OK within a single DB/single run, but
	// if the server restarts, the mappings will get duplicated.
	// So, we need to fix the seed, and think about when/how it
	// is changed. This could be a value stored in the DB?
	fixedSeed := viper.GetInt("hash.seed")
	seed := *(*maphash.Seed)(unsafe.Pointer(&fixedSeed))
	s.h.SetSeed(seed)

	return nil
}

// Init() will be called repeatedly during a single run; specifically,
// before each flush of the buffers. Therefore, everything here should
// be safe to do over-and-over during the life of the service.
func (s *SqliteStorage) Init() error {

	var name string

	if s.Path == ":memory:" {
		name = s.Path
	} else {
		if strings.Contains(s.Filename, ".sqlite") {
			// If they gave us a filename, we might be doing an ETL or similar, and loading a
			// specific, pre-existing database.
			name = fmt.Sprintf("%s/%s?_time_format=sqlite", s.Path, s.Filename)
		} else {
			// If not, we're probably in the middle of logging, and should just open the next database
			// with the appropriate name, since we log by-day at the moment.
			// Heads up: we need to set the date format, because SQLite is really just treating dates as text.
			t := time.Now()
			name = fmt.Sprintf("%s/%s.sqlite?_time_format=sqlite", s.Path, t.Format("2006-01-02"))
		}
	}

	return _init(s, name)
}

// FIXME: replace these with the filename override, above
func (s *SqliteStorage) InitByDate(date time.Time) error {
	name := fmt.Sprintf("%s/%s.sqlite?_time_format=sqlite", s.Path, date.Format("2006-01-02"))
	return _init(s, name)
}

func (s *SqliteStorage) InitByName(filename string) error {
	// FIXME: All of these want some sanitization; no ../asdf.sqlite, etc.
	name := fmt.Sprintf("%s/%s.sqlite?_time_format=sqlite", s.Path, filename)
	return _init(s, name)
}

// This is a massive hole in the abstraction.
// It is used for random data generation, so that I can
// reach down to the SQLC queries directly, and generate
// random events with random timestamps.
func (s *SqliteStorage) GetQueries() *models.Queries {
	return s.queries
}

func hashValue(hash maphash.Hash, s string) int64 {
	if s == "" {
		return 0
	}
	hash.Write([]byte(s))
	h := hash.Sum64()
	hash.Reset()
	return int64(h)
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
			source_h := hashValue(s.h, e.Source)
			event_h := hashValue(s.h, e.Event)

			cluster_h := hashValue(s.h, e.Cluster)
			valid_cluster := false
			if cluster_h != 0 {
				valid_cluster = true
			}

			value_h := hashValue(s.h, e.Value)
			valid_value := false
			if value_h != 0 {
				valid_value = true
			}

			_, err := qtx.LogClusteredEventWithValue(context.Background(), models.LogClusteredEventWithValueParams{
				Timestamp:   e.Timestamp,
				ClusterHash: sql.NullInt64{Int64: cluster_h, Valid: valid_cluster},
				SourceHash:  source_h,
				EventHash:   event_h,
				ValueHash:   sql.NullInt64{Int64: value_h, Valid: valid_value},
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
				Timestamp:  e.Timestamp,
				SourceName: e.Source,
				EventName:  e.Event,
				SourceHash: source_h,
				EventHash:  event_h,
			})
			if err != nil {
				log.Println("Error in storing dictionary:" + err.Error())
				return -1, err
			}

			if valid_value {
				qtx.UpdateLookup(ctx, models.UpdateLookupParams{
					Timestamp: e.Timestamp,
					Name:      e.Value,
					Hash:      value_h,
				})
				if err != nil {
					log.Println("Error in storing value lookup:" + err.Error())
					return -1, err
				}
			}

			counter += 1
		}
	}

	tx.Commit()

	return counter, nil
}

func (s *SqliteStorage) Close() {
	// NOTE: If we're using an in-memory DB, we should not close the DB.
	// This will erase it. Also note, memory DBs are only used for testing
	// at this time. There's probably a better way...
	if !strings.Contains(s.Path, ":memory:") {
		s.db.Close()
	}
}

// -1 if there was an error, 0 if the row was not found, 1 if it was found
func (s *SqliteStorage) TestEventExists(source string, event string) int64 {
	// First hash the value
	source_h := hashValue(s.h, source)
	event_h := hashValue(s.h, event)

	// Check if it can be found in the events table.
	res, err := s.queries.TestEventPairExists(context.Background(), models.TestEventPairExistsParams{
		SourceHash: source_h,
		EventHash:  event_h,
	})
	// If there was an error, just return false.
	if err != nil {
		log.Println(err)
		return -1
	}
	// If it wasn't found, return an error now.
	if res != 1 {
		return res
	}

	// Now do the same with the dictionary
	res, err = s.queries.TestDictionaryPairExists(context.Background(), models.TestDictionaryPairExistsParams{
		SourceHash: source_h,
		EventHash:  event_h,
	})
	// If there was an error, just return false.
	if err != nil {
		return -1
	}
	return res
}

func (s *SqliteStorage) GetDB() *sql.DB {
	return s.db
}

// func (s *SqliteStorage) Event(e *itslog.Event) (int64, error) {
// 	cluster_h := hashValue(s.h, e.Cluster)
// 	source_h := hashValue(s.h, e.Source)
// 	event_h := hashValue(s.h, e.Event)
// 	value_h := hashValue(s.h, e.Value)

// 	valid_cluster := false
// 	valid_value := false
// 	if cluster_h != 0 {
// 		valid_cluster = true
// 	}
// 	if value_h != 0 {
// 		valid_value = true
// 	}

// 	// This is an unsigned to signed conversion...
// 	id, err := s.queries.LogClusteredEventWithValue(context.Background(), models.LogClusteredEventWithValueParams{
// 		ClusterHash: sql.NullInt64{Int64: cluster_h, Valid: valid_cluster},
// 		SourceHash:  source_h,
// 		EventHash:   event_h,
// 		ValueHash:   sql.NullInt64{Int64: value_h, Valid: valid_value},
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	return id, nil
// }
