package schema

import (
	_ "embed"
)

//go:embed schema.sql
var DDL string
