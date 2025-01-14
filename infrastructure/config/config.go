package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL            string
	AutoReconnect  bool
	MaxReconnects  int
	ReconnectDelay int
}

type LoggingConfig struct {
	Level      string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var appConfig *Config

// ReadConfig populates configurations from environment variables.
func Init() {

	appConfig = &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.file", "app.log")
	viper.SetDefault("logging.maxSize", 10)
	viper.SetDefault("logging.maxBackups", 5)
	viper.SetDefault("logging.maxAge", 30)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Failed to unmarshal file: %s", err)
	}
}

// Get private instance config
func Get() *Config {
	return appConfig
}
