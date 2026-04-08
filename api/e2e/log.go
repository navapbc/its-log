package e2e

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand/v2"

	"time"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/spf13/viper"
)

var _versions = [...]string{"v1", "v2", "v3"}
var _endpoints = [...]string{"eob", "patient", "claim"}

func buildBase() string {
	return fmt.Sprintf("http://%s:%s", viper.GetString("serve.host"), viper.GetString("serve.port"))
}

func GenerateLogEvents(iterations int, jitter int, date string) int {
	// This is convoluted, but makes sure we exercise both endpoints.
	var targetUrl string
	if time.Now().Format("2006-01-02") == date {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE
	} else {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	}
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")
	var counter int = 0
	for _, version := range _versions {
		for _, endpoint := range _endpoints {
			// This generates 45 events
			for range iterations {
				time.Sleep(time.Duration(rand.IntN(jitter)) * time.Millisecond)
				bundle := map[string]any{
					"tags": [...]string{version, endpoint},
					"date": date,
				}
				post(targetUrl, bundle, apiKey)
				counter += 1
			}
		}
	}
	return counter
}

func generatePatientId(length int) string {
	bytes := make([]byte, length)
	if _, err := crand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GenerateLogEventsWithValues(iterations int, jitter int, date string) int {
	var targetUrl string
	if time.Now().Format("2006-01-02") == date {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE
	} else {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	}
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range _versions {
		for _, endpoint := range _endpoints {
			// This generates 9*N events
			for range iterations {
				time.Sleep(time.Duration(rand.IntN(jitter)) * time.Millisecond)
				bundle := map[string]any{
					"tags":  [...]string{version, endpoint},
					"value": generatePatientId(8),
					"date":  date,
				}
				post(targetUrl, bundle, apiKey)
				counter += 1
			}
		}
	}
	return counter
}

// 9 events total in three clusters of three
func GenerateClusteredLogs(iterations int, jitter int, date string) int {
	var targetUrl string
	if time.Now().Format("2006-01-02") == date {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE
	} else {
		targetUrl = buildBase() + "/v1" + constants.LOG_CREATE_DATE
	}
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range _versions {
		for _ = range iterations {
			time.Sleep(time.Duration(rand.IntN(jitter)) * time.Millisecond)
			cluster := generatePatientId(8)
			for _, endpoint := range _endpoints {
				bundle := map[string]any{
					"cluster": cluster,
					"tags":    [...]string{version, endpoint},
					"value":   generatePatientId(8),
				}
				post(targetUrl, bundle, apiKey)
				counter += 1
			}
		}
	}
	return counter
}
