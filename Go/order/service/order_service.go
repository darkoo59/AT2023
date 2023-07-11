package service

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type OrderService struct {
	db *gorm.DB
}

func CreateOrderService() *OrderService {
	loadConfig, err := config.LoadConfig()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	databaseURL := loadConfig.DatabaseURL
	db, err := gorm.Open(postgres.Open(databaseURL+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&model.Order{})
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return &OrderService{db}
}

func (service *OrderService) Insert(order *model.Order) (id string, err error) {
	createdId := service.db.Create(order)
	if createdId.Error != nil {
		log.Panic(createdId.Error.Error())
		return "", createdId.Error
	}

	return fmt.Sprint(createdId), nil
}

func (service *OrderService) GetById(id string) (*model.Order, error) {
	var order model.Order
	if err := service.db.First(&order, "id = ?", id).Error; err != nil {
		log.Panic(err.Error())
		return nil, err
	}

	return &order, nil
}
