package b2

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/navapbc/its-log/internal/schema/models"
)

//go:embed etl/sql
var defaultSql embed.FS

//go:embed etl/sequence
var defaultSeq embed.FS

// //go:embed etl/golang
// var defaultGolang embed.FS

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func insertEtl(s *Storage, key, name, kind, body string) {
	err := s.Queries.InsertETL(context.Background(), models.InsertETLParams{
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

func loadFilesFromFS(s *Storage, dirName string) {
	var filesystem embed.FS
	switch dirName {
	case "sql":
		filesystem = defaultSql
	case "sequence":
		filesystem = defaultSeq
	case "golang":
		// filesystem = defaultGolang
	}
	dirEntries, err := fs.ReadDir(filesystem, filepath.Join("etl", dirName))
	if err != nil {
		panic("cannot read embedded directory: " + dirName)
	}

	// Load the default SQL for ETL
	for _, entry := range dirEntries {
		filename := entry.Name()
		if !strings.HasPrefix(filename, "internal_") {
			asBytes, err := fs.ReadFile(filesystem, filepath.Join("etl", dirName, filename))
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

func LoadDefaultEtlFiles(s *Storage) {
	// Just log if we can't read from the embedded FS.
	// Actually... panic? Yeah. This shouldn't fail.
	loadFilesFromFS(s, "sql")
	loadFilesFromFS(s, "sequence")
	//loadFilesFromFS(s, "golang")
}
