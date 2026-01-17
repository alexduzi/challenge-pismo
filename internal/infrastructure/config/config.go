package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	Port    string
	GinMode string
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
		port = "8080" // Default fallback
	}

	viper.SetDefault("GIN_MODE", "debug") // debug, release, or test

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	config := &Config{
		Port:    port,
		GinMode: viper.GetString("GIN_MODE"),
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
