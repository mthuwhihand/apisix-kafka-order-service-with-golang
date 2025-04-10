package main

import (
	"net/http"

	"github.com/apache/apisix-go-plugin-runner/pkg/log"
)

func main() {
	startSSEServer()
	go consumeCreatedOrders("kafka:29092", "created_orders")
	runPlugin()
}

func startSSEServer() {
	http.Handle("/events", sseManager) // path để client connect nhận message
	go func() {
		log.Infof("SSE server running on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Errorf("SSE server error: %v", err)
		}
	}()
}
