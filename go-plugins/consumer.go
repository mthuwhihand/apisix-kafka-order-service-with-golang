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
			clientID := string(msg.Key)
			if clientID == "" {
				log.Warnf("Message key (clientID) is empty for message: %s", string(msg.Value))
				continue
			}

			log.Infof("Sending to clientID %s: %s", clientID, string(msg.Value))
			sseManager.SendToClient(clientID, string(msg.Value))
		} else {
			log.Errorf("Consumer error: %v (%v)", err, msg)
		}
	}
}
