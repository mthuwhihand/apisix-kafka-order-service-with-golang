# Build stage (Bullseye base)
FROM golang:1.23.0-bullseye AS builder

WORKDIR /build

COPY . .

RUN apt-get update && apt-get install -y librdkafka-dev
RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -o go-runner

# Runner image (same glibc version)
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates librdkafka-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /runner

COPY --from=builder /build/go-runner .

RUN chmod +x ./go-runner

CMD ["./go-runner"]
