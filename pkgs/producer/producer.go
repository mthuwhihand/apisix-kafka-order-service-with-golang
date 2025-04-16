package producer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	logger = log.New(log.Writer(), "[producer/producer.go] ", log.LstdFlags|log.Lshortfile)
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(broker, topic string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	kp := &KafkaProducer{
		Producer: producer,
		Topic:    topic,
	}

	go kp.handleEvents()
	return kp, nil
}

func (kp *KafkaProducer) handleEvents() {
	for e := range kp.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				logger.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			} else {
				logger.Printf("Produced event to topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}
		}
	}
}

// func (kp *KafkaProducer) SendMessage(key string, value string) error {
// 	return kp.Producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
// 		Key:            []byte(key),
// 		Value:          []byte(value),
// 	}, nil)
// }

func (kp *KafkaProducer) SendMessage(key string, status_code int, message string, value []byte) error {
	return kp.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		Headers: []kafka.Header{
			{
				Key:   "status_code",
				Value: []byte(strconv.Itoa(status_code)),
			},
			{
				Key:   "message",
				Value: []byte(message),
			},
			{
				Key:   "source",
				Value: []byte("apisix-plugin"),
			},
		},
	}, nil)
}

func (kp *KafkaProducer) Close() {
	kp.Producer.Flush(15000)
	kp.Producer.Close()
}
