package base

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/maphash"
	"log"
	"os"
	"path"
	"time"
	"unsafe"

	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/base/models"
	"github.com/spf13/viper"

	_ "modernc.org/sqlite"
)

func (s *SqliteStorage) ValidateSqliteStorage() error {
	err := validate.Struct(s)

	if err != nil {
		messages := make([]string, 0)
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				messages = append(messages,
					fmt.Sprintf("Field: %s, Tag: %s",
						fieldError.Field(), fieldError.Tag()))

			}
		}

		return fmt.Errorf("storage validation: %s", messages)
	}
	return nil
}

// Init() will be called repeatedly during a single run; specifically,
// before each flush of the buffers. Therefore, everything here should
// be safe to do over-and-over during the life of the service.
func (s *SqliteStorage) Init() error {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())
	}

	switch s.Kind {
	case InMemory:
		log.Println("initializing in-memory SQLite database")
		err := s.ValidateSqliteStorage()
		if err != nil {
			return err
		}
		return s.initMemory()
	case CurrentDatabase:
		// This is a database "for today," e.g. "appid_2026-01-01.sqlite" on 2026/01/01
		t := time.Now()
		s.Filename = fmt.Sprintf("%s_%s.sqlite?_time_format=sqlite", s.AppId, t.Format("2006-01-02"))
		s.Basename = fmt.Sprintf("%s_%s.sqlite", s.AppId, t.Format("2006-01-02"))
	case NamedDatabase:
		// This is a database we choose the name of. For "internal" use only.
		// Always construct the DB name from the key values
		s.Filename = fmt.Sprintf("%s_%s.sqlite?_time_format=sqlite", s.AppId, s.Date)
		s.Basename = fmt.Sprintf("%s_%s.sqlite", s.AppId, s.Date)
	}

	err := s.ValidateSqliteStorage()
	if err != nil {
		return err
	}
	return s.init()
}

func (s *SqliteStorage) Close() {
	// NOTE: If we're using an in-memory DB, we should not close the DB.
	// This will erase it. Because in-memory is only used for testing, this
	// is safe. We are unlikely to use in-memory for production, hence this
	// behavior makes sense at this time.
	if s.Kind != InMemory {
		s.db.Close()
	}
}

// This pulls a constant seed and dupes the maphash
// library into using it every run as the same seed.
// This seed matters for how strings are mapped to integers
// If it changes from one run to the next, then the mapping
// will change. That is OK within a single DB/single run, but
// if the server restarts, the mappings will get duplicated.
// So, we need to fix the seed, and think about when/how it
// is changed. This could be a value stored in the DB?
func (s *SqliteStorage) FixedSeed() {
	fixedSeed := viper.GetInt("hash.seed")
	seed := *(*maphash.Seed)(unsafe.Pointer(&fixedSeed))
	s.h.SetSeed(seed)

}

func (s *SqliteStorage) initMemory() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		return err
	}
	s.db = db

	// create tables
	if _, err := db.ExecContext(ctx, DDL); err != nil {
		return err
	}

	s.queries = models.New(db)
	s.FixedSeed()
	return nil
}

func exists(path string) bool {
	fileInfo, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && fileInfo.Size() > 100
}

func (s *SqliteStorage) init() error {
	db_path := path.Join(s.Path, s.Basename)
	fileExists := exists(db_path)

	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		return err
	}
	s.db = db

	s.queries = models.New(db)
	s.FixedSeed()

	if !fileExists {
		// Create the tables.
		if _, err := db.ExecContext(context.Background(), DDL); err != nil {
			return err
		}
		// If the file exists, check that there's something in the ETL.
		summaries, err := s.queries.GetAllSummaries(context.Background())
		if err != nil {
			return err
		}
		if !(len(summaries) > 0) {
			// If we can't find the table, it isn't initialized.
			log.Println("Loading default ETL values")

			s.LoadDefaultEtlFiles()
		}
	}

	return nil
}
