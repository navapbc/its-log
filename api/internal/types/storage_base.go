package types

import (
	"database/sql"
	"sync"

	"github.com/navapbc/its-log/internal/schema"
)

type Storage struct {
	AppId    string
	Path     []string
	ILTime   *ILTime
	Filename string
	firstUse bool
	lock     sync.Mutex
	db       *sql.DB
	Queries  *schema.Queries
}
