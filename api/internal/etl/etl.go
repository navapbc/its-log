package etl

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/navapbc/its-log/internal/schema"
	"github.com/navapbc/its-log/internal/types"
)

//go:embed sql
var defaultSql embed.FS

//go:embed sequence
var defaultSeq embed.FS

//go:embed golang
var defaultGolang embed.FS

//go:embed starlark
var defaultStarlark embed.FS

func LoadEtlSQLAsString(filename string) string {
	sql, err := defaultSql.ReadFile(path.Join("sql", filename))
	if err != nil {
		panic(err)
	}
	return string(sql)
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func insertEtl(s *types.Storage, key, name, kind, body string) {
	err := s.Queries.InsertETL(context.Background(), schema.InsertETLParams{
		KeyID: key,
		Name:  name,
		Kind:  kind,
		// A golang entry will not have a body.
		Body: sql.NullString{String: body, Valid: true},
	})
	if err != nil {
		log.Printf("could not store SQL in ETL table: %s, %s\n", s.AppId, name)
		log.Printf("err: %s", err.Error())
	}
}

func loadFilesFromFS(s *types.Storage, dirName string) {
	var filesystem embed.FS
	switch dirName {
	case "sql":
		filesystem = defaultSql
	case "sequence":
		filesystem = defaultSeq
	case "golang":
		filesystem = defaultGolang
	case "starlark":
		filesystem = defaultStarlark
	}
	dirEntries, err := fs.ReadDir(filesystem, filepath.Join(dirName))
	if err != nil {
		panic("cannot read embedded directory: " + dirName)
	}

	// Load the default SQL for ETL
	for _, entry := range dirEntries {
		filename := entry.Name()
		if !strings.HasPrefix(filename, "internal_") {
			asBytes, err := fs.ReadFile(filesystem, filepath.Join(dirName, filename))
			if err != nil {
				log.Printf("unable to read file: %v\n", err)
				panic("failed to read file from embedded FS")
			}
			asString := string(asBytes)
			entryName := fileNameWithoutExtension(entry.Name())
			insertEtl(s, "its-log", entryName, dirName, asString)
		}
	}

}

func LoadDefaultEtlFiles(s *types.Storage, from string) error {

	// If the file exists, check that there's something in the ETL.
	_, err := s.Queries.GetETL(context.Background(), "sentinel")
	// If there is nothing in the table, we need to load the ETL
	// into the DB. It is a fresh DB.
	if err != nil {
		// DEBUG LOG
		// log.Printf("LoadDefaultEtlFiles from: %s\n", from)

		// These files are embedded in the app's
		// magic filesystem. If we can't read these,
		// we panic. It shouldn't happen.
		loadFilesFromFS(s, "sql")
		loadFilesFromFS(s, "sequence")
		loadFilesFromFS(s, "golang")
		loadFilesFromFS(s, "starlark")
	}

	return nil
}
