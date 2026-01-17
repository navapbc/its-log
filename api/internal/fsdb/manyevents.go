package fsdb

import (
	"context"
	"database/sql"
	"log"

	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
)

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
				KeyID:       e.KeyId,
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
				KeyID:      e.KeyId,
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
				err = qtx.UpdateLookup(ctx, models.UpdateLookupParams{
					Timestamp: e.Timestamp,
					KeyID:     e.KeyId,
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
