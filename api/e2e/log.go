package e2e

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/navapbc/its-log/internal/base"
)

var versions = [...]string{"v1", "v2", "v3"}
var endpoints = [...]string{"eob", "patient", "claim"}

func buildBase() string {
	return fmt.Sprintf("http://%s:%s", os.Getenv("ITSLOG_SERVE_HOST"), os.Getenv("ITSLOG_SERVE_PORT"))
}

func GenerateLogEvents() int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")
	counter := 0
	for _, version := range versions {
		for _, endpoint := range endpoints {
			// This generates 45 events
			for range 5 {
				bundle := map[string]any{
					"tags": [...]string{version, endpoint},
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
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GenerateLogEventsWithValues() int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range versions {
		for _, endpoint := range endpoints {
			// This generates 45 events
			for range 5 {
				bundle := map[string]any{
					"tags":  [...]string{version, endpoint},
					"value": generatePatientId(8),
				}
				post(targetUrl, bundle, apiKey)
				counter += 1
			}
		}
	}
	return counter
}

// 9 events total in three clusters of three
func GenerateClusteredLogs() int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range versions {
		cluster := generatePatientId(8)
		for _, endpoint := range endpoints {
			bundle := map[string]any{
				"cluster": cluster,
				"tags":    [...]string{version, endpoint},
				"value":   generatePatientId(8),
			}
			post(targetUrl, bundle, apiKey)
			counter += 1
		}
	}
	return counter
}
