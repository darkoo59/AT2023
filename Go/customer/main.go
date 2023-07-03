package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to create logger", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(logger)
	sugar := logger.Sugar()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9000"
	}

	router := gin.Default()
	router.Use(MiddlewareContentTypeSet())

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	sugar.Infof("Server listening on port: %s", port)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))
	//routerGroup := router.Group("/api")

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			sugar.Fatal(err)
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
		logger.Fatal("Cannot gracefully shutdown...")
	}
	sugar.Error("Server stopped")
}
