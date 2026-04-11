package cmd

import (
	"log"
	"os"

	"github.com/navapbc/its-log/endpoints"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/objectstore"
	"github.com/navapbc/its-log/internal/types"
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
	err := base.ConfirmEnvVars()
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	// Load the API keys from the environment
	err = base.GetApiKeys()
	log.Printf("found %d keys", len(base.LiveKeys))
	if err != nil {
		log.Println(err.Error())
		os.Exit(-2)
	}

	log.Printf("storage path: %s", viper.GetString("storage.path"))

	os := objectstore.NewObjectStore("aws")
	os.SetBucket(viper.GetString("backup.bucket"))
	os.SetEndpoint(
		viper.GetString("backup.aws.endpoint.scheme"),
		viper.GetString("backup.aws.endpoint.host"),
		viper.GetString("backup.aws.endpoint.port"))
	os.SetRegion(viper.GetString("backup.aws.region"))
	os.SetKeyIdAndAccessKey(viper.GetString("backup.aws.keyid"), viper.GetString("backup.aws.accesskey"))
	os.Init()
	os.Write([]byte("hello"))

	// FIXME: this should not default to debug
	endpoints.Serve(types.ServeParams{
		Mode: viper.GetString("ginmode"),
	})

}

func init() {
	rootCmd.AddCommand(serveCmd)
}
