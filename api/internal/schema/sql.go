package schema

import (
	"embed"
	_ "embed"
)

//go:embed sql/schema.sql
var DDL string

//go:embed sql
var SQLDir embed.FS
