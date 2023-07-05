package service

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type CustomerService struct {
	logger *zap.Logger
	db     *gorm.DB
}

func CreateCustomerService(logger *zap.Logger) *CustomerService {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal(err.Error())
		return nil
	}

	databaseURL := loadConfig.DatabaseURL
	db, err := gorm.Open(postgres.Open(databaseURL+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&model.Customer{})
	if err != nil {
		logger.Fatal(err.Error())
		return nil
	}

	return &CustomerService{logger, db}
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
