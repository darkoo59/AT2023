package controller

import (
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomerController struct {
	logger  *zap.Logger
	service *service.CustomerService
}

func NewUserController(logger *zap.Logger, service *service.CustomerService) *CustomerController {
	return &CustomerController{logger: logger, service: service}
}

func (controller *CustomerController) CustomerRoute(rg *gin.RouterGroup) {
}
