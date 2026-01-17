package fsdb

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed sql
var defaultSql embed.FS

func (s *SqliteStorage) LoadDefaultEtlSql() {
	dirEntries, err := fs.ReadDir(defaultSql, "sql")

	if err != nil {
		panic("cannot read embedded SQL directory")
	}

	for _, entry := range dirEntries {
		log.Println(entry)
	}
}
