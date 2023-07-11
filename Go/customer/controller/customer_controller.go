package controller

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/middleware"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CustomerController struct {
	logger  *zap.Logger
	service *service.CustomerService
}

func NewUserController(logger *zap.Logger, service *service.CustomerService) *CustomerController {
	return &CustomerController{logger: logger, service: service}
}

func (controller *CustomerController) CustomerRoute(rg *gin.RouterGroup) {
	router := rg.Group("/customer")
	router.POST("/order", middleware.DeserializeCustomer(controller.service, controller.logger),
		controller.Order)
	router.GET("/items", middleware.DeserializeCustomer(controller.service, controller.logger),
		controller.GetItems)
	router.GET("/balance", middleware.DeserializeCustomer(controller.service, controller.logger),
		controller.GetBalance)
	router.PATCH("/balance", middleware.DeserializeCustomer(controller.service, controller.logger),
		controller.UpdateBalance)
}

func (controller *CustomerController) getCustomerIdFromContext(ctx *gin.Context) (string, bool) {
	customer, exists := ctx.Value("currentCustomer").(*model.Customer)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
		return "", true
	}
	customerId := fmt.Sprint(customer.ID)
	return customerId, false
}

type UpdateBalanceRequest struct {
	Balance float64 `json:"balance"`
}

func (controller *CustomerController) Order(ctx *gin.Context) {
	customerId, notLoggedIn := controller.getCustomerIdFromContext(ctx)
	if notLoggedIn {
		return
	}

	var orderBody *service.OrderBody
	if err := ctx.ShouldBindJSON(&orderBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
		return
	}

	err := controller.service.Order(orderBody, customerId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"status": "ordering process started"})
}

func (controller *CustomerController) GetItems(ctx *gin.Context) {
	items := controller.service.GetItemsFromMongoDatabase()
	if items == nil {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (controller *CustomerController) GetBalance(ctx *gin.Context) {
	customerId, notLoggedIn := controller.getCustomerIdFromContext(ctx)
	if notLoggedIn {
		return
	}

	balance, err := controller.service.GetBalanceByCustomerId(customerId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, balance)
}

func (controller *CustomerController) UpdateBalance(ctx *gin.Context) {
	customerId, notLoggedIn := controller.getCustomerIdFromContext(ctx)
	if notLoggedIn {
		return
	}

	var newBalance *UpdateBalanceRequest

	if err := ctx.ShouldBindJSON(&newBalance); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
		return
	}

	if err := controller.service.UpdateBalanceByCustomerId(customerId, newBalance.Balance); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
