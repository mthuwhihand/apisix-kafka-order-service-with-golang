package main

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaMessage struct {
	StatusCode string      `json:"status_code"`
	Message    string      `json:"message"`
	Value      interface{} `json:"value"`
}

func RegisterNewConsumer(broker, topic, group string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	log.Printf("Started consuming from topic %s", topic)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			clientID := string(msg.Key)
			if clientID == "" {
				log.Printf("Message missing clientID: %s", string(msg.Value))
				continue
			}

			// Extract headers
			var statusCode, message string
			for _, header := range msg.Headers {
				switch header.Key {
				case "status_code":
					statusCode = string(header.Value) // Extract status code
				case "message":
					message = string(header.Value) // Extract message
				}
			}

			// Prepare the message object
			// kafkaValueString := string(msg.Value) // Convert byte slice to string

			// If the Kafka value is in JSON format, unmarshal it to a map
			var value map[string]interface{}
			err := json.Unmarshal(msg.Value, &value)
			if err != nil {
				log.Printf("Failed to unmarshal Kafka value: %v", err)
				continue
			}

			// Now, you have the key-value pairs inside the 'value' map
			kafkaMsg := KafkaMessage{
				StatusCode: statusCode,
				Message:    message,
				Value:      value, // Assuming msg.Value is the data you want to send
			}

			// Convert the message to JSON
			kafkaMsgJSON, err := json.Marshal(kafkaMsg)
			if err != nil {
				log.Printf("Failed to marshal Kafka message: %v", err)
				continue
			}

			// Log and send to client
			log.Printf("Sending to client %s: %s", clientID, string(kafkaMsgJSON))

			// Send to the client with the JSON message
			sseManager.SendToClient(clientID, string(kafkaMsgJSON))
		} else {
			log.Printf("Consumer error: %v", err)
		}
	}
}
