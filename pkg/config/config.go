package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort          string        `mapstructure:"SERVER_PORT"`
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBUrl               string        `mapstructure:"DB_URL"`
	RedisAddr           string        `mapstructure:"REDIS_ADDR"`
	RedisPassword       string        `mapstructure:"REDIS_PASSWORD"`
	JWTSecret           string        `mapstructure:"JWT_SECRET"`
	JWTExpirationHours  int           `mapstructure:"JWT_EXPIRATION_HOURS"`
	StripeSecretKey     string        `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret string        `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	AccessTokenDuration time.Duration `mapstructure:"-"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	config.AccessTokenDuration = time.Duration(config.JWTExpirationHours) * time.Hour

	return
}
