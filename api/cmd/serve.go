package cmd

import (
	"log"
	"os"

	"github.com/navapbc/its-log/endpoints"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	err := base.ConfirmEnvVars()
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	defer base.ConfigureLoggers(viper.GetString("ginmode")).Sync()

	// Load the API keys from the environment
	err = base.GetApiKeys()
	zap.L().Debug("found API keys", zap.Int("length", len(base.LiveKeys)))
	if err != nil {
		zap.L().Error(err.Error())
		os.Exit(-2)
	}

	zap.L().Info("storage path", zap.String("storage.path", viper.GetString("storage.path")))

	// FIXME: Need CTRL-C/OS signal handling
	endpoints.Serve(types.ServeParams{
		Mode: viper.GetString("ginmode"),
	})

}

func init() {
	rootCmd.AddCommand(serveCmd)
}
