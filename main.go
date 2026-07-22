// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// // TelemetryPayload defines the structure of incoming telemetry logs
// type TelemetryPayload struct {
// 	DeviceID  string  `json:"device_id"`
// 	Metric    string  `json:"metric"`
// 	Value     float64 `json:"value"`
// 	Timestamp int64   `json:"timestamp"`
// }

// func telemetryHandler(w http.ResponseWriter, r *http.Request) {
// 	// Restrict endpoint to POST requests only
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var payload TelemetryPayload
// 	// Decode JSON request body into the struct
// 	err := json.NewDecoder(r.Body).Decode(&payload)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Default timestamp if missing
// 	if payload.Timestamp == 0 {
// 		payload.Timestamp = time.Now().Unix()
// 	}

// 	// Print incoming telemetry log to terminal for validation
// 	fmt.Printf("[RECEIVED] Device: %s | Metric: %s | Value: %.2f | Time: %d\n",
// 		payload.DeviceID, payload.Metric, payload.Value, payload.Timestamp)

// 	// Send 200 OK Response
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"status": "accepted"}`))
// }

// func main() {
// 	http.HandleFunc("/api/v1/telemetry", telemetryHandler)

// 	fmt.Println("🚀 Telemetry API Gateway running on port 8080...")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Printf("Server failed to start: %v\n", err)
// 	}
// }

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type TelemetryData struct {
	DeviceID string  `json:"device_id"`
	Metric   string  `json:"metric"`
	Value    float64 `json:"value"`
	Time     int64   `json:"timestamp,omitempty"`
}

var kafkaWriter *kafka.Writer

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "telemetry-events",
		Balancer: &kafka.LeastBytes{},
	}
}

func telemetryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data TelemetryData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if data.Time == 0 {
		data.Time = time.Now().Unix()
	}

	payload, _ := json.Marshal(data)

	// Send event directly to Kafka
	err := kafkaWriter.WriteMessages(r.Context(), kafka.Message{
		Key:   []byte(data.DeviceID),
		Value: payload,
	})

	if err != nil {
		log.Printf("Failed to produce to Kafka: %v", err)
		http.Error(w, "Failed to publish telemetry", http.StatusInternalServerError)
		return
	}

	log.Printf("[KAFKA PRODUCED] Device: %s | Metric: %s | Value: %.2f", data.DeviceID, data.Metric, data.Value)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}

func main() {
	initKafka()
	defer kafkaWriter.Close()

	http.HandleFunc("/api/v1/telemetry", telemetryHandler)

	fmt.Println("🚀 Telemetry API Gateway running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}