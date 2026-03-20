package base

import (
	_ "embed"

	"github.com/go-playground/validator/v10"
)

// https://pkg.go.dev/github.com/go-playground/validator/v10
var validate *validator.Validate

//go:embed schema.sql
var DDL string

// type SqliteType int

// const (
// 	InMemory SqliteType = iota + 1
// 	CurrentDatabase
// 	NamedDatabase
// )

// func (s SqliteStorage) Validate() error {
// 	switch s.Kind {
// 	case InMemory:
// 		return nil
// 	case CurrentDatabase:
// 		if len(s.Path) < 1 {
// 			return fmt.Errorf("len(Path): %d", len(s.Path))
// 		}
// 	case NamedDatabase:
// 		if len(s.Path) < 1 || len(s.Filename) < 1 {
// 			return fmt.Errorf("len(Path): %d, len(Filename): %d", len(s.Path), len(s.Filename))
// 		}
// 	}
// 	return nil
// }

// // https://pkg.go.dev/github.com/go-playground/validator/v10
// type SqliteStorage struct {
// 	Kind     SqliteType `validate:"required"`
// 	Path     string     `validate:"required_unless=Kind 1"`
// 	Filename string     `validate:"required_unless=Kind 1"`
// 	Basename string
// 	Date     string `validate:"required"`
// 	// Keep separate users separate with this value
// 	// It will come in via the environment
// 	AppId string `validate:"required"`
// 	// sqlc queries
// 	queries *models.Queries
// 	// For hashing str to int, consistently
// 	h maphash.Hash
// 	// For rare cases where we have to reach through sqlc
// 	db *sql.DB
// }
