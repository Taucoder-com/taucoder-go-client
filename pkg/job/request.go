package job

import (
	"encoding/json"
	"log"
	"net/http"
)

func DoRequest(req *http.Request) (*JobStatusResponse, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var status JobStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return &status, nil
}
