package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

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

func (s *SqliteStorage) Event(source string, event string, value any, value_type string) (int64, error) {
	fmt.Printf("%s %v %v\n", event, value, value_type)

	id, err := s.queries.LogIt(context.Background(), models.LogItParams{
		Version: "v1",
		Source:  source,
		Event:   event,
		Value:   fmt.Sprintf("%v", value),
		Type:    value_type,
	})

	if err != nil {
		panic(err)
	}

	return id, nil
}
