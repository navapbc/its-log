package types

import (
	"database/sql"
	"time"
)

type Event struct {
	Timestamp time.Time
	AppId     string
	KeyId     string
	Cluster   string   `json:"cluster"`
	Tags      []string `json:"tags"`
	TagString string
	Value     string `json:"value"`
	Date      string `json:"date"`
}

type EventRow struct {
	ID        string
	Timestamp string
	KeyId     string
	Cluster   sql.NullString
	Tags      string
	Value     sql.NullString
}
