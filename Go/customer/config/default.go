package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Port                   string        `mapstructure:"PORT"`
	DatabaseURL            string        `mapstructure:"DATABASE_URL"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRES_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRES_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	ActorCustomerAddress   string        `mapstructure:"ACTOR_CUSTOMER_ADDRESS"`
	ActorCustomerPort      int           `mapstructure:"ACTOR_CUSTOMER_PORT"`
	ActorOrderAddress      string        `mapstructure:"ACTOR_ORDER_ADDRESS"`
	ActorOrderPort         int           `mapstructure:"ACTOR_ORDER_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to create logger", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		logger.Info(err.Error())
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	return config, err
}
