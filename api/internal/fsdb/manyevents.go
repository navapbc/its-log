package fsdb

import (
	"context"
	"database/sql"
	"log"

	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
)

/*
	tags_h := hashValue(s.h, e.TagString)
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

	key_h := hashValue(s.h, e.KeyId)
*/

func (s *SqliteStorage) ManyEvents(evt_buff []*itslog.Event) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	counter := int64(0)
	qtx := s.queries.WithTx(tx)
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

			// 20260226 MCJ
			// An earlier design involved turning all of these values into hashes using
			// the function hashValue(s.h, e.Value) or similar. We're simplifying the
			// code to record strings. We're still storing a hash map in the lookup table,
			// as the cost is negligible and it leaves us a path back if we want the data
			// compression savings from using INT64s for all our strings instead of... strings.
			_, err := qtx.LogEvent(context.Background(), models.LogEventParams{
				Timestamp: e.Timestamp,
				KeyID:     e.KeyId,
				Cluster:   sql.NullString{String: e.Cluster, Valid: valid_cluster},
				Tags:      e.TagString,
				Value:     sql.NullString{String: e.Value, Valid: valid_value},
			})

			if err != nil {
				log.Println("Error in storing event:" + err.Error())
				return -1, err
			}

			// Use the transaction to update the dictionary
			// in bulk as well. Individual inserts should
			// quietly ignore conflicts. This could be optimized to only update
			// when we see a new hash value.
			tags_h := hashValue(s.h, e.TagString)
			err = qtx.UpdateLookup(ctx, models.UpdateLookupParams{
				Timestamp: e.Timestamp,
				Kind:      "tags",
				Hash:      tags_h,
				Name:      e.TagString,
			})
			if err != nil {
				log.Println("Error in tags in lookup:" + err.Error())
				return -1, err
			}

			if valid_value {
				value_h := hashValue(s.h, e.Value)
				err = qtx.UpdateLookup(ctx, models.UpdateLookupParams{
					Timestamp: e.Timestamp,
					Name:      e.Value,
					Kind:      "value",
					Hash:      value_h,
				})
				if err != nil {
					log.Println("Error in storing value lookup:" + err.Error())
					return -1, err
				}
			}

			if valid_cluster {
				cluster_h := hashValue(s.h, e.Cluster)
				err = qtx.UpdateLookup(ctx, models.UpdateLookupParams{
					Timestamp: e.Timestamp,
					Name:      e.Cluster,
					Kind:      "cluster",
					Hash:      cluster_h,
				})
				if err != nil {
					log.Println("Error in storing cluster lookup:" + err.Error())
					return -1, err
				}
			}

			counter += 1
		}
	}

	err = tx.Commit()
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
