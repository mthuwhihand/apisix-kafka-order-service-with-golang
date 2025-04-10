package main

import (
	_ "hihand/docs"
	config "hihand/internal/configs/dev"
	"hihand/internal/router"
	"log"
)

var (
	logger = log.New(log.Writer(), "[server/main.go] ", log.LstdFlags|log.Lshortfile)
)

// @title Hihand API
// @version 1.0
// @description API For Order Service
// @host localhost:8080
// @BasePath /
func main() {
	config, cfgErr := config.Instance()
	if cfgErr != nil {
		logger.Println("Can not get config:", cfgErr)
	}

	// Khởi tạo Kafka Producer và Consumer
	broker := config.BROKER
	topic := config.TOPIC_ORDER
	groupID := "order_consumer_group"
	responseTopic := config.TOPIC_ORDER_CREATED

	producer, kafkaConsumer, err := router.StartOrderKafkaConsumer(broker, topic, groupID, responseTopic)
	if err != nil {
		log.Fatalf("Error initializing Kafka: %v", err)
	}
	defer producer.Close()
	defer kafkaConsumer.Close()

	app := router.NewRouter()

	app.Run(":8080")

	logger.Println("App Running!")

}
