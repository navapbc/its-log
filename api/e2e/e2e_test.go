package e2e

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/navapbc/its-log/internal/base"
	"github.com/spf13/viper"
)

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

func configViper() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// Everything we want in the env for Viper
	// must be prefixed with ITSLOG_
	viper.SetEnvPrefix("ITSLOG")
	viper.AutomaticEnv()
}

func TestE2E(t *testing.T) {
	configViper()
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

	RunTests()
}
