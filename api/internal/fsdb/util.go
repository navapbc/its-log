package fsdb

import (
	"fmt"
	"hash/maphash"
	"path"

	"github.com/jadudm/its-log/internal/fsdb/models"
)

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

// The goal is to make sure we have *only* a filename.
func MakeGoodSqliteFilename(fname string) (string, error) {
	cleaned := path.Base(fname)
	// We could end up with "." or "/", so make sure the length is greater
	// than one, and that it ends in ".sqlite"
	if len(cleaned) > 1 {
		return cleaned, nil
	}
	return "BADSQLITEFILENAME", fmt.Errorf("given an unrecoverable SQLite filename: %s", fname)
}
