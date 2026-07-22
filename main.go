package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TelemetryPayload defines the structure of incoming telemetry logs
type TelemetryPayload struct {
	DeviceID  string  `json:"device_id"`
	Metric    string  `json:"metric"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

func telemetryHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict endpoint to POST requests only
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload TelemetryPayload
	// Decode JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Default timestamp if missing
	if payload.Timestamp == 0 {
		payload.Timestamp = time.Now().Unix()
	}

	// Print incoming telemetry log to terminal for validation
	fmt.Printf("[RECEIVED] Device: %s | Metric: %s | Value: %.2f | Time: %d\n",
		payload.DeviceID, payload.Metric, payload.Value, payload.Timestamp)

	// Send 200 OK Response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "accepted"}`))
}

func main() {
	http.HandleFunc("/api/v1/telemetry", telemetryHandler)

	fmt.Println("🚀 Telemetry API Gateway running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}