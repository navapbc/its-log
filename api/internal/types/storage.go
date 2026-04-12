package types

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/navapbc/its-log/internal/schema"
	"github.com/navapbc/its-log/internal/schema/models"

	// _ "github.com/ncruces/go-sqlite3/driver"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

// with ncruces
// https://github.com/ncruces/go-sqlite3/commit/7c820ede3caf7a861f53c700e188db59d9928d3b
// var SQLITE_PARAMS = [...]string{"_timefmt=sqlite"}
// with modernc
var SQLITE_PARAMS = [...]string{"_time_format=sqlite"}

// with ncruces
// const SQLITE_DRIVER string = "sqlite3"
// with modernc
const SQLITE_DRIVER string = "sqlite"

type Storage struct {
	AppId    string
	Path     []string
	ILTime   *ILTime
	Filename string
	firstUse bool
	lock     sync.Mutex
	db       *sql.DB
	Queries  *models.Queries
}

// We name things based on dates
func (s *Storage) YYYYMMDD() string {
	return s.ILTime.AsYYYYMMDD()
}

func (s *Storage) SetDateYMD(ymd string) error {
	if ymd == "today" {
		s.ILTime = NewILTimeToday()
		return nil
	}
	d, e := ILTimeFromYMD(ymd)
	if e == nil {
		s.ILTime = d
	}
	return e
}

func (s *Storage) SetDateILT(ilt *ILTime) {
	s.ILTime = ilt
}

func NewStorage(appId string) *Storage {
	// Create DB
	s := &Storage{
		AppId: appId,
	}
	s.ILTime = NewILTimeToday()

	return s
}

func (s *Storage) SubtractDays(days int) {
	s.ILTime.SubtractDays(days)
}

func (s *Storage) AddDays(days int) {
	s.ILTime.AddDays(days)
}

func (s *Storage) Init() error {
	// Compute filename based on date and app ID
	s.Filename = fmt.Sprintf("%s_%s.sqlite", s.AppId, s.YYYYMMDD())
	s.Path = []string{viper.GetString("storage.path"), s.Filename}

	dbParams := ""
	if len(SQLITE_PARAMS) > 0 {
		dbParams = "?" + strings.Join(SQLITE_PARAMS[:], "&")
	}
	dbPath := path.Join(path.Join(s.Path...) + dbParams)
	// Create tables
	log.Printf("Storage.Init: opening %s\n", s.Filename)
	db, err := sql.Open(SQLITE_DRIVER, "file:"+dbPath)
	if err != nil {
		panic("could not create database: " + s.Filename)
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
	// log.Println("locking")
	s.lock.Lock()
}

func (s *Storage) Unlock() {
	// log.Println("unlocking")
	s.lock.Unlock()
}

func (s *Storage) Delete() {
	path := path.Join(viper.GetString("storage.path"), s.Filename)
	os.Remove(path)
}
