package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	Port     string
	GinMode  string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("PORT")
	}
	if port == "" {
		port = "8080"
	}

	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("APP_NAME", "api-accounts")

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "accounts_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_CONNS", 25)
	viper.SetDefault("DB_MIN_CONNS", 5)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	config := &Config{
		AppName: viper.GetString("APP_NAME"),
		Port:    port,
		GinMode: viper.GetString("GIN_MODE"),
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			MaxConns: viper.GetInt("DB_MAX_CONNS"),
			MinConns: viper.GetInt("DB_MIN_CONNS"),
		},
	}

	if config.Database.Host == "" || config.Database.User == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	AppConfig = config
	return config, nil
}

func GetConfig() *Config {
	if AppConfig == nil {
		log.Fatal("Config not initialized. Call LoadConfig() first.")
	}
	return AppConfig
}
