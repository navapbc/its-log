package etl

import (
	"github.com/navapbc/its-log/internal/types"
)

// Any file that starts with internal_ will not be pulled
// into the SQLite table. The map is not actual ETL code,
// so we leave it out of the table.
var GolangETLMap = map[string]func(types.EtlParams){
	"count-all-combinations": CountAllCombinations,
}
