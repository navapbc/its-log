package types

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/duckdb/duckdb-go/v2"
	"github.com/navapbc/its-log/internal/schema"
	"github.com/spf13/viper"
)

func NewStorage(appId string) (*Storage, error) {
	// Create DB
	s := &Storage{
		AppId: appId,
	}
	s.ILTime = NewILTimeToday()
	s.db = nil

	return s, nil
}

func (s *Storage) InitDB() error {
	// Compute filename based on date and app ID
	s.Filename = fmt.Sprintf("%s_%s.duckdb", s.AppId, s.YYYYMMDD())
	s.Path = []string{viper.GetString("storage.path"), s.Filename}

	// Don't re-initialize a database
	if s.db != nil {
		return nil
	}

	dbPath := path.Join(path.Join(s.Path...))
	// Create tables
	// DEBUG LOG
	// log.Printf("Storage.Init: opening %s\n", s.Filename)
	c, err := duckdb.NewConnector(dbPath, func(execer driver.ExecerContext) error {
		// https://duckdb.org/docs/current/configuration/overview
		bootQueries := []string{
			`SET threads TO 1`,
			// `SET schema=main`,
			// `SET search_path=main`,
		}
		for _, query := range bootQueries {
			_, err := execer.ExecContext(context.Background(), query, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	db := sql.OpenDB(c)

	s.db = db

	// Load query models
	s.Queries = schema.New(s.db)

	err = s.InitTables()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) InitTables() error {
	// Create tables
	_, err := s.db.ExecContext(context.Background(), schema.DDL)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func (s *Storage) Lock() {
	// log.Println("locking")
	s.lock.Lock()
}

func (s *Storage) Unlock() {
	// log.Println("unlocking")
	s.lock.Unlock()
}

func (s *Storage) Delete() {
	s.db.Close()
	path := path.Join(viper.GetString("storage.path"), s.Filename)
	log.Printf("Deleting database: %s", path)
	os.Remove(path)
}
