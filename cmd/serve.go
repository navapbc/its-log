/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jadudm/its-log/internal/config"
	defaultstorage "github.com/jadudm/its-log/internal/default-storage"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/serve"
	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the its-log API",
	Long: `Serve the its-log API
This takes no parameters; all configureation is provided
via the application's 'config.yaml'.
`,
	Run: serve_cmd,
}

func serve_cmd(cmd *cobra.Command, args []string) {
	// main() should set everything up before
	// handing it off to the engine to execute.

	// Instantiate the preferred storage backend
	// This may become a choice at some point.
	var storage itslog.ItsLog
	switch viper.GetString("app.storage") {
	case "sqlite":
		storage = &sqlite.SqliteStorage{
			Path: viper.GetString("sqlite.path"),
		}
		err := storage.Init()
		if err != nil {
			panic(err)
		}
	case "default":
		storage = &defaultstorage.DefaultStorage{}
		storage.Init()
	}

	// Parse the API key config
	var apiConfig config.ApiKeys
	err := viper.Unmarshal(&apiConfig)
	// TODO: Handle config failure
	if err != nil {
		panic(err)
	}

	serve.Serve(storage, apiConfig)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
