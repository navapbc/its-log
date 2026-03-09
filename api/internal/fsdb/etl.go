package fsdb

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/jadudm/its-log/internal/fsdb/models"
)

//go:embed sql
var defaultSql embed.FS

func (s *SqliteStorage) GetDB() *sql.DB {
	return s.db
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (s *SqliteStorage) LoadDefaultEtlSql() {
	dirEntries, err := fs.ReadDir(defaultSql, "sql")

	if err != nil {
		panic("cannot read embedded SQL directory")
	}

	for _, entry := range dirEntries {
		filename := entry.Name()
		sqlAsBytes, err := fs.ReadFile(defaultSql, filepath.Join("sql", filename))
		if err != nil {
			log.Printf("unable to read file: %v\n", err)
			return
		}
		sqlAsString := string(sqlAsBytes)
		sqlName := fileNameWithoutExtension(entry.Name())
		err = s.queries.InsertETL(context.Background(), models.InsertETLParams{
			KeyID: 0, // Use a default key ID of 0 for the automatically inserted values
			Name:  sqlName,
			Sql:   sqlAsString,
		})
		if err != nil {
			log.Printf("could not store SQL in ETL table: %s, %s\n", s.AppId, sqlName)
			log.Printf("err: %s", err.Error())
		}
	}
}
