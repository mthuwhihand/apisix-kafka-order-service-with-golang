#!/bin/bash

echo "Reading .env..."
cat /opt/kafka/.env

source /opt/kafka/.env

echo "Creating topics: $TOPIC_ORDER, $TOPIC_ORDER_CREATED"

# Gọi đúng tên CLI tool được cài sẵn
kafka-topics --bootstrap-server kafka:29092 --create --if-not-exists --topic "$TOPIC_ORDER" --partitions 1 --replication-factor 1

kafka-topics --bootstrap-server kafka:29092 --create --if-not-exists --topic "$TOPIC_ORDER_CREATED" --partitions 1 --replication-factor 1

echo "Kafka topics created successfully!"
