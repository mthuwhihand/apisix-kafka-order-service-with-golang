package controllers

import (
	"hihand/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	logger = log.New(log.Writer(), "[server/main.go] ", log.LstdFlags|log.Lshortfile)
)

type OrderController struct {
	service services.OrderService
}

func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{service: service}
}

// HelloWorld godoc
// @Summary Hello world
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /orders/hello-world [get]
func (c *OrderController) HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Hello world"})
}

// CreateOrder godoc
// @Summary Create New Order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order info"
// @Example {"product_id": 123, "quantity": 2, "customer_name": "John Doe", "address": "123 Main Street"}
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	// var order models.Order
	// if err := ctx.ShouldBindJSON(&order); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// err := c.service.Create(&order)
	// if err != nil {
	// 	logger.Println("Failed to create order:", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
	// 	return
	// }

	// ctx.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": order})

	ctx.JSON(200, map[string]string{
		"message": "Message pushed to Kafka",
	})
	ctx.Abort()
}

// UpdateOrder godoc
// @Summary Update Order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param updates body map[string]interface{} true "Update fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [patch]
func (c *OrderController) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	var updates map[string]interface{}

	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Update(orderID, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder godoc
// @Summary Delete Order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [delete]
func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	err := c.service.Delete(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (c *OrderController) SearchOrders(ctx *gin.Context) {
	query := ctx.Query("query")
	limit := ctx.DefaultQuery("limit", "10")
	skip := ctx.DefaultQuery("skip", "0")

	limitInt, _ := strconv.Atoi(limit)
	skipInt, _ := strconv.Atoi(skip)

	orders, err := c.service.Search(query, limitInt, skipInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search orders"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}
