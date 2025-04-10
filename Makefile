hello:
	echo "Hello"
# Define the run-server target
run-server:
	go run ./cmd/server/main.go

# Define the run-go-plugins
run-go-plugins:
	go run ./apisix/go-plugins/main.go

tidy:
	go mod tidy


#1/ Run go plugins
#APISIX_LISTEN_ADDRESS=unix:/tmp/runner.sock ./apisix/go-plugins/apisix-go-plugin run
#2/ Run docker compose up --build
#docker compose -f docker-compose.yml up --build
#docker compose -f kafka.docker-compose.yml up --build
#docker compose -f apisix.docker-compose.yml up --build
#3/ make run-server

#Check logs:
#docker exec -it kafka /usr/bin/kafka-topics --bootstrap-server localhost:9092 --list
# Optional: Add a default target
.PHONY: all
all: run-server