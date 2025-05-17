# apisix-kafka-order-service-with-golang

A microservices-based project using Golang, Kafka, and APISIX to handle order processing via event-driven architecture.

## 🧰 Technologies Used

* Golang
* Kafka
* APISIX
* Docker Compose
* NodeJS

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/mthuwhihand/apisix-kafka-order-service-with-golang.git
cd apisix-kafka-order-service-with-golang
```

### 2. Start services

For start db
```bash
docker compose -f docker-compose.yml up
```

For start kafka
```bash
docker compose -f kafka.docker-compose.yml up
```

For start local Backend service
```bash
make run-server
```

For start APISIX Api Gateway 
```bash
docker compose -f apisix.docker-compose.yml up
```

For start local SSE (Server Send Event) service
```bash
cd sse-service && go run .
```

For start FE
```bash
cd hihand-fe && npm run dev
```


## 📁 Project Structure

* `cmd/server/` – Entry point for the main API service
* `internal/` – Business logic, Kafka handling, etc.
* `socket-service/`, `sse-service/` – Real-time communication services
* `hihand-fe/` – Frontend
* `docker-compose.yml` – Container orchestration
---
