package e2e

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

var _versions = [...]string{"v1", "v2", "v3"}
var _endpoints = [...]string{"eob", "patient", "claim"}

func buildBase() string {
	return fmt.Sprintf("http://%s:%s", viper.GetString("serve.host"), viper.GetString("serve.port"))
}

func checkRes(res map[string]any) {
	if v, ok := res["status"].(string); ok {
		if v != "ok" {
			log.Println("POST ERROR: " + res["message"].(string))
		}
	}
}

func wiggle(jitter int) {
	// time.Sleep(time.Duration(rand.IntN(jitter)) * time.Millisecond)
	// Make this more deterministic
	time.Sleep(time.Duration(jitter * int(time.Millisecond)))
}

func GenerateLogEvents(iterations int, jitter int, date *types.ILTime) int {
	// This is convoluted, but makes sure we exercise both endpoints.
	var targetUrl string
	targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")
	var counter int = 0
	for _, version := range _versions {
		for _, endpoint := range _endpoints {
			// This generates 45 events
			for range iterations {
				wiggle(jitter)
				bundle := map[string]any{
					"tags": [...]string{version, endpoint},
					"date": date.AsYYYYMMDD(),
				}

				res := post(targetUrl, bundle, apiKey)
				checkRes(res)

				counter += 1
			}
		}
	}

	// log.Printf("GenerateLogEvents: %d\n", counter)

	return counter
}

func generatePatientId(length int) string {
	bytes := make([]byte, length)
	if _, err := crand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GenerateLogEventsWithValues(iterations int, jitter int, date *types.ILTime) int {
	var targetUrl string

	targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range _versions {
		for _, endpoint := range _endpoints {
			// This generates 9*N events
			for range iterations {
				wiggle(jitter)
				bundle := map[string]any{
					"tags":  [...]string{version, endpoint},
					"value": generatePatientId(8),
					"date":  date.AsYYYYMMDD(),
				}

				res := post(targetUrl, bundle, apiKey)
				checkRes(res)

				counter += 1
			}
		}
	}

	return counter
}

// 9 events total in three clusters of three
func GenerateClusteredLogs(t *testing.T, iterations int, jitter int, date *types.ILTime) int {
	var targetUrl string
	targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range _versions {
		for _ = range iterations {
			wiggle(jitter)
			cluster := generatePatientId(8)
			for _, endpoint := range _endpoints {
				bundle := map[string]any{
					"cluster": cluster,
					"tags":    [...]string{version, endpoint},
					"value":   generatePatientId(8),
					"date":    date.AsYYYYMMDD(),
				}

				res := post(targetUrl, bundle, apiKey)
				checkRes(res)

				counter += 1
			}
		}
	}
	return counter
}
