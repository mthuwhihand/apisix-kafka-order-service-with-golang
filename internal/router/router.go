package router

import (
	"encoding/json"
	"fmt"
	"hihand/internal/controllers"
	"hihand/internal/database"
	"hihand/internal/models"
	"hihand/internal/services"
	"hihand/pkgs/consumer"
	"hihand/pkgs/producer"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	logger = log.New(log.Writer(), "[router/router.go] ", log.LstdFlags|log.Lshortfile)
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	logger.Println("database.Instance()")
	db, err := database.Instance()
	if err != nil {
		logger.Println("Get Database Instance failed! Error:", err)
	}

	logger.Println("database.AutoMigrate()")
	err = database.AutoMigrate()
	if err != nil {
		logger.Println("Migration failed:", err)
	} else {
		logger.Println("Database migration completed successfully!")
	}

	// Middleware
	r.Use(gin.Logger())   //Log request
	r.Use(gin.Recovery()) //Catch panic or something like that, and return 500 Internal Server Error
	service := services.NewOrderService(db)
	controller := controllers.NewOrderController(service)

	api := r.Group("/orders")
	{
		api.GET("", controller.SearchOrders)
		api.POST("", controller.CreateOrder)
		api.PATCH("", controller.UpdateOrder)
		api.DELETE("/:id", controller.DeleteOrder)
	}
	r.GET("/hello-world", controller.HelloWorld)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	return r
}

func StartOrderKafkaConsumer(broker, topic, groupID, responseTopic string) (*producer.KafkaProducer, *consumer.KafkaConsumer, error) {
	db, err := database.Instance()
	if err != nil {
		logger.Println("Get Database Instance failed! Error:", err)
	}
	service := services.NewOrderService(db)

	producer, err := producer.NewKafkaProducer(broker, responseTopic)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return nil, nil, err
	}

	processMessage := func(msg *kafka.Message) ([]byte, error) {
		logger.Printf("Processing message: %s\n", string(msg.Value))

		var order models.Order

		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			logger.Printf("Failed to parse message: %v\n", err)
			return nil, fmt.Errorf("failed to parse message: %v", err)
		}

		err = service.Create(&order)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %v", err)
		}

		jsonBytes, err := json.Marshal(order)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value to JSON: %v", err)
		}

		return jsonBytes, nil
	}

	kafkaConsumer, err := consumer.NewKafkaConsumer(broker, topic, groupID, responseTopic, producer, processMessage)
	if err != nil {
		log.Printf("Failed to create Kafka consumer: %v", err)
		return nil, nil, err
	}

	log.Printf("New Kafka consumer created successfully. Listening on broker: %s, topic: %s, groupID: %s\n", broker, topic, groupID)

	return producer, kafkaConsumer, nil
}
