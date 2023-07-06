package service

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/customer_actor"
	messages "github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/message"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type CustomerService struct {
	logger        *zap.Logger
	db            *gorm.DB
	customerActor *customer_actor.CustomerActor
}

func CreateCustomerService(customerActor *customer_actor.CustomerActor, zapLogger *zap.Logger) *CustomerService {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		zapLogger.Fatal(err.Error())
		return nil
	}

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

	err = db.AutoMigrate(&model.Customer{})
	if err != nil {
		zapLogger.Fatal(err.Error())
		return nil
	}

	return &CustomerService{zapLogger, db, customerActor}
}

type OrderBody struct {
	ItemId   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

func (service *CustomerService) Insert(customer *model.Customer) (id string, err error) {
	createdId := service.db.Create(customer)
	if createdId.Error != nil {
		service.logger.Error(createdId.Error.Error())
		return "", createdId.Error
	}

	return fmt.Sprint(createdId), nil
}

func (service *CustomerService) GetByEmail(email string) (*model.Customer, error) {
	var customer model.Customer
	if err := service.db.First(&customer, "email = ?", email).Error; err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return &customer, nil
}

func (service *CustomerService) GetById(id string) (*model.Customer, error) {
	var customer model.Customer
	if err := service.db.First(&customer, "id = ?", id).Error; err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return &customer, nil
}

func (service *CustomerService) EmailExists(email string) bool {
	var customer model.Customer
	if err := service.db.First(&customer, "email = ?", email).Error; err != nil {
		return false
	}

	return true
}

func (service *CustomerService) Order(order *messages.ReceiveOrder_Request) error {
	pid := service.customerActor.Spawn()
	service.customerActor.Send(pid, order)
	service.logger.Info("Message sent to customer-actor")
	return nil
}
