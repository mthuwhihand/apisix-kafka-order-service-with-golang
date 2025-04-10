package main

import (
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var sseManager = NewSSEManager()

func consumeCreatedOrders(broker, topic string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "created_order_consumer_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Errorf("Failed to create consumer: %v", err)
		return
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	log.Infof("Started consuming from topic %s", topic)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Infof("Received message from %s: %s", topic, string(msg.Value))
			sseManager.Broadcast(string(msg.Value))
		} else {
			log.Errorf("Consumer error: %v (%v)", err, msg)
		}
	}
}
