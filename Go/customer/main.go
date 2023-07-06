package main

import (
	"context"
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/controller"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/customer_actor"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/service"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func MiddlewareContentTypeSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func initializeControllers(routerGroup *gin.RouterGroup, customerActor *customer_actor.CustomerActor, logger *zap.Logger) {
	customerService := service.CreateCustomerService(customerActor, logger)
	customerController := controller.NewUserController(logger, customerService)
	customerController.CustomerRoute(routerGroup)

	authService := service.CreateAuthService(logger, customerService)
	authController := controller.NewAuthController(logger, authService)
	authController.AuthRoute(routerGroup)
}

func initializeActorSystem(loadConfig config.Config, logger *zap.Logger) *customer_actor.CustomerActor {
	system := actor.NewActorSystem()
	remoteConfig := remote.Configure(loadConfig.ActorHostAddress, loadConfig.ActorCustomerPort)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	actorContext := system.Root
	customerActor := &customer_actor.CustomerActor{
		Remoting: remoting,
		Context:  actorContext,
		Logger:   logger,
	}
	customerActorProps := actor.PropsFromProducer(func() actor.Actor {
		return customerActor
	})
	remoting.Register("customer-actor", customerActorProps)
	logger.Info("Customer actor registered")

	return customerActor
}

func main() {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	logger, err := zapConfig.Build()
	if err != nil {
		fmt.Println("Failed to create logger", err.Error())
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(logger)
	sugar := logger.Sugar()

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		sugar.Error(err.Error())
	}

	port := loadConfig.Port
	if len(port) == 0 {
		port = "9000"
	}
	router := gin.Default()
	router.Use(MiddlewareContentTypeSet())
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	sugar.Infof("Server listening on port: %s", port)

	corsConfig := corsConfig()
	router.Use(cors.New(corsConfig))
	routerGroup := router.Group("/api")

	customerActor := initializeActorSystem(loadConfig, logger)
	initializeControllers(routerGroup, customerActor, logger)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			sugar.Info(err.Error())
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	sugar.Infof("Received terminate, graceful shutdown [%s]", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if server.Shutdown(timeoutContext) != nil {
		logger.Error("Cannot gracefully shutdown...")
	}
	sugar.Info("Server stopped")
}

func corsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	return corsConfig
}
