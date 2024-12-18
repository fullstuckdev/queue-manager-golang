package config

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	RabbitMQ struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Server struct {
		Port string
	}
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	
	cfg := &Config{}
	
	cfg.DB.Host = viper.GetString("DB_HOST")
	cfg.DB.Port = viper.GetString("DB_PORT")
	cfg.DB.User = viper.GetString("DB_USER")
	cfg.DB.Password = viper.GetString("DB_PASSWORD")
	cfg.DB.Name = viper.GetString("DB_NAME")

	cfg.RabbitMQ.Host = viper.GetString("RABBITMQ_HOST")
	cfg.RabbitMQ.Port = viper.GetString("RABBITMQ_PORT")
	cfg.RabbitMQ.User = viper.GetString("RABBITMQ_USER")
	cfg.RabbitMQ.Password = viper.GetString("RABBITMQ_PASSWORD")

	cfg.Server.Port = viper.GetString("SERVER_PORT")

	return cfg, nil
} 