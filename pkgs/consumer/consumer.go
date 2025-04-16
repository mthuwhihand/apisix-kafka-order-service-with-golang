package consumer

import (
	"fmt"
	"hihand/pkgs/producer"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	logger = log.New(log.Writer(), "[consumer/consumer.go] ", log.LstdFlags|log.Lshortfile)
)

type KafkaConsumer struct {
	Consumer         *kafka.Consumer
	Producer         *producer.KafkaProducer // Producer là tùy chọn
	Topic            string
	ResponseTopic    string                                   // ResponseTopic là tùy chọn
	ProcessMessageFn func(msg *kafka.Message) ([]byte, error) // Callback function xử lý message
}

func NewKafkaConsumer(broker, topic, groupID, responseTopic string, producer *producer.KafkaProducer, processMessageFn func(msg *kafka.Message) ([]byte, error)) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"sasl.mechanisms":   "PLAIN",
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %v", topic, err)
	}

	kc := &KafkaConsumer{
		Consumer:         consumer,
		Producer:         producer,
		Topic:            topic,
		ResponseTopic:    responseTopic,
		ProcessMessageFn: processMessageFn,
	}

	go kc.handleEvents()
	return kc, nil
}

func (kc *KafkaConsumer) handleEvents() {
	for {
		msg, err := kc.Consumer.ReadMessage(-1)
		if err == nil {
			logger.Printf("Received message: %s\n", string(msg.Value))

			// Use callback function to process message
			result, err := kc.ProcessMessageFn(msg)
			if err != nil {
				logger.Printf("Error processing message: %v\n", err)
			}

			if kc.Producer != nil && kc.ResponseTopic != "" {
				if err != nil {
					err = kc.Producer.SendMessage(string(msg.Key), http.StatusBadRequest, "Order failed", []byte(err.Error()))
				} else {
					err = kc.Producer.SendMessage(string(msg.Key), http.StatusOK, "Order successfully", result)
				}

				if err != nil {
					logger.Printf("Failed to send message: %v\n", err)
				} else {
					logger.Printf("Successfully sent processed message back to topic %s\nMessage: %s", kc.ResponseTopic, result)
				}
			}
		} else {
			logger.Printf("Error while consuming message: %v\n", err)
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.Consumer.Close()
}
