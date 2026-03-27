package e2e

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand/v2"

	"os"
	"time"

	"github.com/navapbc/its-log/internal/base"
)

var versions = [...]string{"v1", "v2", "v3"}
var endpoints = [...]string{"eob", "patient", "claim"}

func buildBase() string {
	return fmt.Sprintf("http://%s:%s", os.Getenv("ITSLOG_SERVE_HOST"), os.Getenv("ITSLOG_SERVE_PORT"))
}

func GenerateLogEvents(N int64) int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")
	counter := 0
	for _, version := range versions {
		for _, endpoint := range endpoints {
			// This generates 45 events
			for range N {
				time.Sleep(time.Duration(rand.IntN(3)) * time.Millisecond)
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
	if _, err := crand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GenerateLogEventsWithValues(N int64) int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range versions {
		for _, endpoint := range endpoints {
			// This generates 9*N events
			for range N {
				time.Sleep(time.Duration(rand.IntN(3)) * time.Millisecond)
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
func GenerateClusteredLogs(N int64) int {
	targetUrl := buildBase() + "/v1/log"
	apiKey, _ := base.GetKeyBundle("pupper", "pup_logging")

	counter := 0
	for _, version := range versions {
		for _ = range N {
			time.Sleep(time.Duration(rand.IntN(3)) * time.Millisecond)
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
	}
	return counter
}
