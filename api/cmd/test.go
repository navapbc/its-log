package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

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

func stressTest(wg *sync.WaitGroup, iterations int, jitter int) {
	defer wg.Done()
	gleCount := e2e.GenerateLogEvents(iterations, jitter)
	glevCount := e2e.GenerateLogEventsWithValues(iterations, jitter)
	gclCount := e2e.GenerateClusteredLogs(iterations, jitter)
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

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go stressTest(&wg, 50*i, 10*i)
	}
	wg.Wait()

	log.Println(path.Join(os.Getenv("ITSLOG_STORAGE_PATH")))
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(testCmd)
}
