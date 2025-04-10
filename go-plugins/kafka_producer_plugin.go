package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	"github.com/apache/apisix-go-plugin-runner/pkg/runner"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap/zapcore"
)

// KafkaProducerPlugin là plugin để produce message vào Kafka.
type KafkaProducerPlugin struct {
	Broker string `json:"broker"` // Kafka broker (ví dụ: "kafka:9092")
	Topic  string `json:"topic"`  // Topic để gửi message (ví dụ: "orders")
}

func (kp *KafkaProducerPlugin) Name() string {
	return "kafka-producer"
}

func (kp *KafkaProducerPlugin) ParseConf(conf []byte) (interface{}, error) {
	err := json.Unmarshal(conf, kp)
	return kp, err
}

func (kp *KafkaProducerPlugin) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	log.Infof("KafkaProducerPlugin triggered")
	body, err := r.Body()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to get body: %v", err)))
		log.Errorf("Failed to read request body: %v", err)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty body"))
		log.Infof("Received request with empty body")
		return
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kp.Broker,
	})
	if err != nil {
		log.Errorf("Failed to create producer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Producer creation error: %v", err)))
		return
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event, 1)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
		Key:            []byte("order"),
		Value:          body,
	}, deliveryChan)

	if err != nil {
		log.Errorf("Failed to produce message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to produce message: %v", err)))
		return
	}

	e := <-deliveryChan
	m, ok := e.(*kafka.Message)
	if !ok {
		log.Errorf("Kafka event is not a message: %v", e)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Kafka event type mismatch"))
		return
	}

	if m.TopicPartition.Error != nil {
		log.Errorf("Delivery failed: %v", m.TopicPartition.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Delivery failed: %v", m.TopicPartition.Error)))
	} else {
		log.Infof("Message delivered to topic %s [%d] at offset %v",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order request sent to Kafka"))
	}
	close(deliveryChan)
}

func (kp *KafkaProducerPlugin) ResponseFilter(conf interface{}, w pkgHTTP.Response) {
	// Không xử lý response.
}

func runPlugin() {
	listenAddress := os.Getenv("APISIX_LISTEN_ADDRESS")
	if listenAddress == "" {
		// Set mặc định nếu chưa được set
		listenAddress = "unix:/tmp/runner.sock"
		os.Setenv("APISIX_LISTEN_ADDRESS", listenAddress) // để runner.Run dùng được
	}

	log.Infof("Starting APISIX Go Plugin Runner on %s", listenAddress)

	if err := plugin.RegisterPlugin(&KafkaProducerPlugin{
		Broker: "kafka:29092",
		Topic:  "orders",
	}); err != nil {
		log.Errorf("Error registering plugin: %v", err)
		return
	}

	cfg := runner.RunnerConfig{
		LogLevel:  zapcore.DebugLevel,
		LogOutput: os.Stdout,
	}

	runner.Run(cfg)
}
