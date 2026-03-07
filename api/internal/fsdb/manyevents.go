package fsdb

import (
	"context"
	"database/sql"
	"log"

	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
)

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
