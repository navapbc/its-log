package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/e2e"
	serve "github.com/navapbc/its-log/endpoints"
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
		"ITSLOG_GINMODE":             "production", // "debug" or "production"
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

func testSequence(N int64) {
	gleCount := e2e.GenerateLogEvents(N)
	glevCount := e2e.GenerateLogEventsWithValues(N)
	gclCount := e2e.GenerateClusteredLogs(N)
	total := gleCount + glevCount + gclCount

	e2e.RunDefaultSequence()

	log.Printf("total events: %d\n", total)

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
	gin.DefaultWriter = io.Discard
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	go serve.Serve()

	for _ = range 5 {
		go testSequence(100)
	}

	for _ = range 2 {
		time.Sleep(1 * time.Second)
		log.Print("x")
	}
	e2e.RunDefaultSequence()

	for _ = range 2 {
		time.Sleep(1 * time.Second)
		log.Print("y")
	}
	e2e.RunDefaultSequence()

	for _ = range 5 {
		time.Sleep(1 * time.Second)
		log.Print("z")
	}

	e2e.RunDefaultSequence()

	log.Println(path.Join(os.Getenv("ITSLOG_STORAGE_PATH")))
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(testCmd)
}
