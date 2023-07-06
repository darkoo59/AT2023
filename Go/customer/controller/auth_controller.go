package controller

import (
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/service"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	logger  *zap.Logger
	service *service.AuthService
}

func NewAuthController(logger *zap.Logger, service *service.AuthService) *AuthController {
	return &AuthController{logger: logger, service: service}
}

func (ac *AuthController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/register", ac.Register)
	router.POST("/login", ac.Login)
	router.GET("/refresh", ac.RefreshAccessToken)
	router.GET("/logout", ac.Logout)
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var registerBody *service.RegisterBody
	if err := ctx.ShouldBindJSON(&registerBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
		return
	}

	customer := &model.Customer{
		Email:    registerBody.Email,
		Password: registerBody.Password,
	}

	err := ac.service.Register(customer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "successful registration"})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var credentials service.LoginCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
		return
	}

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		ac.logger.Fatal(err.Error())
		return
	}

	id, err := ac.service.Login(credentials)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := util.CreateToken(loadConfig.AccessTokenExpiresIn, id, loadConfig.AccessTokenPrivateKey)
	if err != nil {
		ac.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}
	refreshToken, err := util.CreateToken(loadConfig.RefreshTokenExpiresIn, id, loadConfig.RefreshTokenPrivateKey)
	if err != nil {
		ac.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, loadConfig.RefreshTokenMaxAge*60,
		"/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "refresh token is missing"})
		return
	}

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		ac.logger.Fatal(err.Error())
		return
	}

	sub, err := util.ValidateToken(cookie, loadConfig.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	id, err := ac.service.DoesUserExist(sub.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		ac.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := util.CreateToken(loadConfig.AccessTokenExpiresIn, id, loadConfig.AccessTokenPrivateKey)
	if err != nil {
		ac.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
