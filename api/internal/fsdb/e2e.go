package fsdb

import (
	"database/sql"
)

// -1 if there was an error, 0 if the row was not found, 1 if it was found
// func (s *SqliteStorage) TestEventExists(source string, event string) int64 {
// 	// First hash the value
// 	source_h := hashValue(s.h, source)
// 	event_h := hashValue(s.h, event)

// 	// Check if it can be found in the events table.
// 	res, err := s.queries.TestEventPairExists(context.Background(), models.TestEventPairExistsParams{
// 		SourceHash: source_h,
// 		EventHash:  event_h,
// 	})
// 	// If there was an error, just return false.
// 	if err != nil {
// 		log.Println(err)
// 		return -1
// 	}
// 	// If it wasn't found, return an error now.
// 	if res != 1 {
// 		return res
// 	}

// 	// Now do the same with the dictionary
// 	res, err = s.queries.TestDictionaryPairExists(context.Background(), models.TestDictionaryPairExistsParams{
// 		SourceHash: source_h,
// 		EventHash:  event_h,
// 	})
// 	// If there was an error, just return false.
// 	if err != nil {
// 		return -1
// 	}
// 	return res
// }

func (s *SqliteStorage) GetDB() *sql.DB {
	return s.db
}
