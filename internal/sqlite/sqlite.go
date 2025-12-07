package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/sqlite/models"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

type SqliteStorage struct {
	Path    string
	queries *models.Queries
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

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	s.queries = models.New(db)

	return nil
}

func (s *SqliteStorage) Event(e *itslog.Event) (int64, error) {
	fmt.Printf("%s %v %v\n", e.Event, e.Value, e.Type)

	id, err := s.queries.LogIt(context.Background(), models.LogItParams{
		Version: e.Version,
		Source:  e.Source,
		Event:   e.Event,
		Value:   fmt.Sprintf("%v", e.Value),
		Type:    e.Type,
	})

	if err != nil {
		panic(err)
	}

	return id, nil
}
