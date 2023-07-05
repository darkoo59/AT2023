package middleware

import (
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/service"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func DeserializeCustomer(cs *service.CustomerService, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
			return
		}

		loadConfig, err := config.LoadConfig(".")
		if err != nil {
			logger.Fatal(err.Error())
			return
		}

		sub, err := util.ValidateToken(accessToken, loadConfig.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		user, err := cs.GetById(sub.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		ctx.Set("currentCustomer", user)
		ctx.Next()
	}
}
