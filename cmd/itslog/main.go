package main

import (
	"fmt"

	"cms.hhs.gov/its-log/internal/itslog"
	"cms.hhs.gov/its-log/internal/sqlite"
	"github.com/spf13/viper"
)

func main() {
	// We must explicitly read the config
	// before doing anything else.
	ReadConfig()
	// I should abstract the storage engine.
	var s itslog.ItsLog

	// Instantiate the preferred storage backend
	switch viper.GetString("storage") {
	case "sqlite":
		s = &sqlite.SqliteStorage{
			Path: viper.GetString("sqlite_path"),
		}
	case "s3":
		// pass
	default:
		// pass
	}

	err := s.Init()
	if err != nil {
		panic(err)
	}

	engine := PourGin(s)

	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}
