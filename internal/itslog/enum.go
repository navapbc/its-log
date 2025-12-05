package itslog

type LogItType int64

const (
	INTEGER LogItType = iota
	REAL
	TEXT
	DATE
	DATETIME
	JSONB
	BLOB
)
