package cmd

import (
	"encoding/json"

	"github.com/jadudm/its-log/internal/config"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/serve"
	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/madflojo/testcerts"
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
	storage = &sqlite.SqliteStorage{
		Path: viper.GetString("storage.path"),
	}

	// Parse the API key config
	var apiConfig config.ApiKeys
	key_string := viper.GetString("apikeys")
	err := json.Unmarshal([]byte(key_string), &apiConfig)
	// TODO: Handle config failure
	if err != nil {
		panic(err)
	}

	// Set up for TLS.
	if viper.GetString("serve.cert") == "mock" && viper.GetString("serve.key") == "mock" {
		// If we end up on an immutable filesystem, this will become an issue.
		cert, key, err := testcerts.GenerateCertsToTempFile("/tmp")
		if err != nil {
			panic(err)
		}
		viper.Set("serve.cert", cert)
		viper.Set("serve.key", key)
		serve.Serve(storage, apiConfig)
	}

	// This will fail if we don't have a cert/key set in the config.
	serve.Serve(storage, apiConfig)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
