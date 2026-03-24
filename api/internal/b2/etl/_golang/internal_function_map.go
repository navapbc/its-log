package etl

import (
	"database/sql"

	"github.com/navapbc/its-log/internal/base"
)

// Any file that starts with internal_ will not be pulled
// into the SQLite table. The map is not actual ETL code,
// so we leave it out of the table.
var GolangETLMap = map[string]func(s *base.ServeCtx, tx *sql.Tx){
	"count-all-combinations": CountAllCombinations,
}
