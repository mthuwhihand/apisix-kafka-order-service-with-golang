package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "hihand/internal/configs/dev"
	"hihand/internal/models"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
	logger     = log.New(log.Writer(), "[databases] ", log.LstdFlags|log.Lshortfile)
)

func Instance() (*gorm.DB, error) {
	var err error

	dbOnce.Do(func() {
		cfg, cfgErr := config.Instance()
		if cfgErr != nil {
			logger.Println("Can not get config:", cfgErr)
			err = cfgErr
			return
		}

		for i := 0; i < 10; i++ { // Retry up to 10 times
			db, dbErr := gorm.Open(postgres.Open(cfg.DB_URL), &gorm.Config{})
			if dbErr != nil {
				logger.Println("Can not connect to database, retrying...:", dbErr)
				time.Sleep(5 * time.Second)
			} else {
				dbInstance = db
				break
			}
		}

		if dbInstance == nil {
			err = fmt.Errorf("failed to connect to database after multiple attempts")
		} else {
			logger.Println("Connected to PostgreSQL successfully!")
		}
	})

	return dbInstance, err
}

func AutoMigrate() error {
	return dbInstance.AutoMigrate(
		&models.Order{},
		&models.OrderDetail{},
	)
}
