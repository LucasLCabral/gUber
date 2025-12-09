package main

import (
	"bytes"
	"encoding/json"
	"guber/shared/contracts"
	"io"
	"log"
	"net/http"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body once and keep it for forwarding
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the request body
	var reqBody previewTripRequest
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	if reqBody.UserID == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}

	// Forward the request to trip-service with the original body
	resp, err := http.Post("http://trip-service:8083/preview", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Error calling trip-service: %v", err)
		http.Error(w, "failed to connect to trip service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Trip service returned status code: %d", resp.StatusCode)
		http.Error(w, "trip service error", resp.StatusCode)
		return
	}

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Printf("Error decoding trip service response: %v", err)
		http.Error(w, "failed to parse JSON data from trip service", http.StatusBadGateway)
		return
	}

	response := contracts.APIResponse{Data: respBody}
	writeJson(w, http.StatusOK, response)
}
