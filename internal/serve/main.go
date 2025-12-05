package serve

import (
	"fmt"

	"cms.hhs.gov/its-log/internal/itslog"
	"cms.hhs.gov/its-log/internal/sqlite"
	"github.com/spf13/viper"
)

func Serve() {
	// I should abstract the storage engine.
	var s itslog.ItsLog

	// Instantiate the preferred storage backend
	switch viper.GetString("app.storage") {
	case "sqlite":
		s = &sqlite.SqliteStorage{
			Path: viper.GetString("sqlite.path"),
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

	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}
