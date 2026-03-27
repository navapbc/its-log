package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/navapbc/its-log/internal/types"
)

func post(targetUrl string, bundle map[string]interface{}, k *types.ApiKey) {
	jsonBytes, _ := json.MarshalIndent(bundle, "", "  ")

	req, err := http.NewRequest("POST", targetUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Println("did not log: " + string(jsonBytes))
	}
	// Because this is test code, running in a test context, we can
	// pull the key from the same place the running server will pull it.
	// This key would not exist in production, etc.
	req.Header.Set("X-Api-Key", k.Key)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	// Ensure the response body is closed to prevent resource leaks
	resp.Body.Close()
}

func get(targetUrl string, k *types.ApiKey) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Println("did not get: " + targetUrl)
	}
	req.Header.Set("X-Api-Key", k.Key)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		contents, _ := io.ReadAll(resp.Body)
		log.Println("status: " + string(contents))
	}
	resp.Body.Close()
}
