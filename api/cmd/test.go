package cmd

import (
	"log"
	"os"
	"path"

	"github.com/navapbc/its-log/e2e"
	"github.com/navapbc/its-log/internal/base"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the its-log API",
	Long: `This runs an all-up test of the API.
It will generate a database file locally.
`,
	Run: test_cmd,
}

func configureEnv() {
	// These are the *required* env vars for its-log.
	var env = map[string]string{
		"ITSLOG_SERVE_HOST":          "localhost",
		"ITSLOG_SERVE_PORT":          "8888",
		"ITSLOG_APIKEY_PUPPERLOG":    "{\"app_id\": \"pupper\", \"key_id\": \"pup_logging\", \"permission\": \"log\", \"key\": \"1234567890123456123456789012345612345678901234561234567890123456\"}",
		"ITSLOG_APIKEY_PUPPERADMIN":  "{\"app_id\": \"pupper\", \"key_id\": \"pup_admin\", \"permission\": \"admin\", \"key\": \"abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh\"}",
		"ITSLOG_BUFFER_FLUSHWAITSEC": "1",
		"ITSLOG_BUFFER_LENGTH":       "100",
		"ITSLOG_GINMODE":             "debug", // "debug" or "production"
		"ITSLOG_PROXIES_TRUSTED":     "TBD",
		"ITSLOG_STORAGE_PATH":        path.Join(os.TempDir()),
	}
	// Load the environment
	for k, v := range env {
		// Allow command-line overrides.
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

func test_cmd(cmd *cobra.Command, args []string) {
	configureEnv()
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

	e2e.RunTests()
}

func init() {
	rootCmd.AddCommand(testCmd)
}
