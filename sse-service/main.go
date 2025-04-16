package main

import (
	"github.com/gin-gonic/gin"
)

var sseManager = NewSSEManager()

func main() {
	go RegisterNewConsumer("localhost:9092", "created_orders", "sse_order_created_group")

	r := gin.Default()
	r.GET("/events/order_created", sseManager.HandleOrderCreatedResponseSSE)

	r.Run(":8083") // Listen on port 8083
}
