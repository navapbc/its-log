package csp

import (
	"os"
	"path/filepath"

	"github.com/jadudm/its-log/internal/sqlite"
	cron "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

const START_OF_THE_DAY = "2 * * * *"

// TODO
// It's unclear if having a completely isolated gofunc
// that opens up the database from yesterday and does some stuff
// is a good design. That said, we'll be writing to a new DB, so
// this won't have any contention at the time it runs.
// This could also be done externally by a command-line cron job.
func Summarize() {
	c := cron.New()
	c.AddFunc(START_OF_THE_DAY,
		func() {
			path := viper.GetString("sqlite.path")
			files, err := os.ReadDir(path)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				storage := &sqlite.SqliteStorage{
					Path: filepath.Join(path, file.Name()),
				}
				//FIXME: Walk the path, and try and process/summarize everything there.
				// The metadata should tell us if it was already done.
				storage.Init()
				storage.Summarize()

			}
		})
	c.Start()
}
