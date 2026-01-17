package cmd

import (
	"log"
	"os"

	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/serve"
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
	// This will panic if we don't have the tools we need.
	err := itslog.ConfirmEnvVars()
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	// Load the API keys from the environment
	err = itslog.GetApiKeys()
	log.Printf("found %d keys", len(itslog.LiveKeys))
	if err != nil {
		log.Println(err.Error())
		os.Exit(-2)
	}

	log.Printf("storage path: %s", viper.GetString("storage.path"))

	serve.Serve()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
