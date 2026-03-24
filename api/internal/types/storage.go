package types

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/navapbc/its-log/internal/schema"
	"github.com/navapbc/its-log/internal/schema/models"
	_ "github.com/ncruces/go-sqlite3/driver"
	"github.com/spf13/viper"
)

var SQLITE_PARAMS = [...]string{"_timefmt=sqlite"}

const SQLITE_DRIVER string = "sqlite3"

type Storage struct {
	AppId    string
	date     time.Time
	filename string
	firstUse bool
	lock     sync.Mutex
	db       *sql.DB
	Queries  *models.Queries
}

// We name things based on dates
func (s *Storage) YYYYMMDD() string {
	return s.date.Format("2006-01-02")
}

func (s *Storage) Today() {
	s.date = time.Now()
}

func (s *Storage) Yesterday() {
	s.date = time.Now().AddDate(0, 0, -1)
}

func (s *Storage) SetDate(ymd string) error {
	d, e := time.Parse("2006-01-02", ymd)
	if e != nil {
		s.date = d
	}
	return e
}

func NewStorage(appId string) *Storage {
	// Create DB
	s := &Storage{
		AppId: appId,
	}
	s.Today()

	return s
}

func (s *Storage) Init() error {
	// Compute filename based on date and app ID
	s.filename = fmt.Sprintf("%s_%s.sqlite", s.AppId, s.YYYYMMDD())
	// https://github.com/ncruces/go-sqlite3/commit/7c820ede3caf7a861f53c700e188db59d9928d3b
	dbPath := path.Join(viper.GetString("storage.path"), s.filename+"?"+strings.Join(SQLITE_PARAMS[:], "&"))
	// Create tables
	db, err := sql.Open(SQLITE_DRIVER, "file:"+dbPath)
	if err != nil {
		panic("could not create database: " + s.filename)
	}
	s.db = db

	// Load query models
	s.Queries = models.New(s.db)

	// Create tables
	_, err = s.db.ExecContext(context.Background(), schema.DDL)
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
	s.lock.Lock()
}

func (s *Storage) Unlock() {
	s.lock.Unlock()
}
