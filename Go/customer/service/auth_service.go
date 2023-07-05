package service

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/util"
	"go.uber.org/zap"
)

type AuthService struct {
	logger          *zap.Logger
	customerService *CustomerService
}

func CreateAuthService(logger *zap.Logger, customerService *CustomerService) *AuthService {
	return &AuthService{logger, customerService}
}

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (as *AuthService) Register(customer *model.Customer) error {
	if !customer.IsEmailValid() {
		return fmt.Errorf("email %s is not valid", customer.Email)
	}
	if as.customerService.EmailExists(customer.Email) {
		return fmt.Errorf("email %s already exists", customer.Email)
	}

	password, err := util.HashPassword(customer.Password)
	if err != nil {
		return err
	}
	customer.Password = password

	_, err = as.customerService.Insert(customer)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) Login(credentials LoginCredentials) (string, error) {
	customer, err := as.customerService.GetByEmail(credentials.Email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}
	if err := util.ComparePasswords(customer.Password, credentials.Password); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	return customer.ID.String(), nil
}

func (as *AuthService) DoesUserExist(id string) (string, error) {
	customer, err := as.customerService.GetById(id)
	if err != nil {
		return "", err
	}

	return customer.ID.String(), nil
}
