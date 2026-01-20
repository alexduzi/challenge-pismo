package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func resetViperAndConfig() {
	viper.Reset()
	AppConfig = nil
}

func TestLoadConfig_WithDefaultValues(t *testing.T) {
	// arrange
	resetViperAndConfig()

	// Garantir que não há variáveis de ambiente configuradas
	os.Unsetenv("PORT")
	os.Unsetenv("GIN_MODE")
	os.Unsetenv("APP_NAME")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("DB_MAX_CONNS")
	os.Unsetenv("DB_MIN_CONNS")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")
	os.Unsetenv("API_VERSION")
	os.Unsetenv("API_TIMEOUT")

	// act
	config, err := LoadConfig()

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "debug", config.GinMode)
	assert.Equal(t, "api-accounts", config.AppName)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "postgres", config.Database.Password)
	assert.Equal(t, "accounts_db", config.Database.DBName)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, 25, config.Database.MaxConns)
	assert.Equal(t, 5, config.Database.MinConns)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "text", config.LogFormat)
	assert.Equal(t, "v1", config.ApiVersion)
	assert.Equal(t, 30, config.ApiTimeout)
	assert.Equal(t, config, AppConfig)
}
