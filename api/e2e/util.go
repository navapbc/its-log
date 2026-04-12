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

func post(targetUrl string, bundle map[string]interface{}, k *types.ApiKey) map[string]any {
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
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	resp.Body.Close()

	asmap := make(map[string]any)
	err = json.Unmarshal(result, &asmap)

	return asmap
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

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
