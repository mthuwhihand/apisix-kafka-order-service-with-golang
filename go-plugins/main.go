package main

import (
	"os"

	"github.com/apache/apisix-go-plugin-runner/pkg/log"
)

func main() {
	listenAddress := os.Getenv("APISIX_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "unix:/tmp/runner.sock"
		os.Setenv("APISIX_LISTEN_ADDRESS", listenAddress)
	}

	runKafkaProducerPlugin()

	log.Infof("Starting APISIX Go Plugin Runner on %s", listenAddress)

}
