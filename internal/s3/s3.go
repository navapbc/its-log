package s3

import (
	_ "embed"
	"fmt"

	"cms.hhs.gov/its-log/internal/itslog"
	"cms.hhs.gov/its-log/internal/sqlite/models"
	_ "modernc.org/sqlite"
)

type SqliteStorage struct {
	Path    string
	Queries *models.Queries
}

func (s *SqliteStorage) Init() error {

	return nil
}

func (s *SqliteStorage) Event(event string, value any, value_type itslog.LogItType) (int64, error) {
	fmt.Printf("%s %v %v\n", event, value, value_type)

	return 0, nil
}
