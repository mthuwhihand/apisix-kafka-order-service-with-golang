package router

import (
	"hihand/internal/controllers"
	"hihand/internal/database"
	"hihand/internal/repositories"
	"hihand/internal/services"
	"log"

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
	repository := repositories.NewOrderRepository(db)
	service := services.NewOrderService(repository)
	controller := controllers.NewOrderController(service)

	api := r.Group("/orders")
	{
		api.POST("", controller.CreateOrder)
		api.PATCH("", controller.UpdateOrder)
		api.DELETE("", controller.DeleteOrder)
	}
	r.GET("/hello-world", controller.HelloWorld)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	return r
}
