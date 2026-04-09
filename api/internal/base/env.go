package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/navapbc/its-log/internal/types"
)

var LiveKeys types.ApiKeys

var minimum = []string{
	"ITSLOG_BUFFER_FLUSHWAITSEC",
	"ITSLOG_BUFFER_LENGTH",
	"ITSLOG_GINMODE",
	"ITSLOG_PROXIES_TRUSTED",
	"ITSLOG_SERVE_HOST",
	"ITSLOG_SERVE_PORT",
	"ITSLOG_STORAGE_PATH",
}

func ConfirmEnvVars() error {
	allKeys := make([]string, 0)
	allValues := make([]string, 0)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		allKeys = append(allKeys, pair[0])
		allValues = append(allValues, pair[1])
	}

	for _, must := range minimum {
		if !slices.Contains(allKeys, must) {
			return fmt.Errorf("missing required environment variable: %s", must)
		}
	}

	return nil
}

func ConfirmApiKeyMapKeys(envvar string, m map[string]string) error {
	keys := slices.Collect(maps.Keys(m))
	must_haves := []string{
		"app_id",
		"key_id",
		"permission",
		"key",
	}

	for _, key := range keys {
		if !slices.Contains(must_haves, key) {
			return fmt.Errorf("%s is not a valid JSON key", key)
		}
	}

	for _, must_have := range must_haves {
		if !slices.Contains(keys, must_have) {
			return fmt.Errorf("%s missing JSON key %s", envvar, must_have)
		}
	}

	return nil
}

// API keys take the form
// ITSLOG_APIKEY_<appId>='{"app_id": "the_app", "key_id": "uniq", "permission": "log", "key": "32-byte-string"}'
func GetApiKeys() error {
	LiveKeys = make([]types.ApiKey, 0)

	allvars := os.Environ()
	for _, e := range allvars {
		key := types.ApiKey{}
		pair := strings.SplitN(e, "=", 2)

		if strings.HasPrefix(pair[0], "ITSLOG_APIKEY_") {
			log.Println("unmarshalling " + pair[0])
			// key.AppId = pieces[2]
			keymap := make(map[string]string)
			err := json.Unmarshal([]byte(pair[1]), &keymap)
			if err != nil {
				return fmt.Errorf("could not parse JSON for %s", pair[0])
			}
			err = ConfirmApiKeyMapKeys(pair[0], keymap)
			if err != nil {
				return err
			}

			key.AppId = keymap["app_id"]
			key.KeyId = keymap["key_id"]
			key.PermissionString = keymap["permission"]
			key.Key = keymap["key"]
			key.ConfigurePermissions()
			LiveKeys = append(LiveKeys, key)
		}
	}

	return nil
}

func GetKeyBundle(appId, keyId string) (*types.ApiKey, error) {
	for _, k := range LiveKeys {
		if k.AppId == appId && k.KeyId == keyId {
			return &k, nil
		}
	}
	return nil, errors.New("could not find key for " + appId + ", " + keyId)
}
