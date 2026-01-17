package fsdb

import (
	"database/sql"
	_ "embed"
	"hash/maphash"

	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/fsdb/models"
)

// https://pkg.go.dev/github.com/go-playground/validator/v10
var validate *validator.Validate

//go:embed schema.sql
var ddl string

type SqliteType int

const (
	InMemory SqliteType = iota + 1
	CurrentDatabase
	NamedDatabase
)

type SqliteStorage struct {
	Kind     SqliteType `validate:"required"`
	Path     string     `validate:"required_if=Kind 3"`
	Filename string     `validate:"required_unless=Kind 1"`
	// Keep separate users separate with this value
	// It will come in via the environment
	AppId string `validate:"required"`
	// sqlc queries
	queries *models.Queries
	// For hashing str to int, consistently
	h maphash.Hash
	// For rare cases where we have to reach through sqlc
	db *sql.DB
}
