package csp

import (
	"time"

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
			yesterday := time.Now().AddDate(0, 0, -1)
			storage := &sqlite.SqliteStorage{
				Path: viper.GetString("sqlite.path"),
			}
			storage.Init(yesterday)
			storage.Summarize()
		})
	c.Start()
}
