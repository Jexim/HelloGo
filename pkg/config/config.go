package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address string
	}
	Database struct {
		URI string
	}
	Logger struct {
		Level string
	}
}

func Load() *Config {
	// Read .env file if present
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	// Allow ENV variables to override, replace "." in keys with "_"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("server.address", ":3000")
	viper.SetDefault("logger.level", "info")

	cfg := &Config{}
	cfg.Server.Address = viper.GetString("server.address")
	cfg.Database.URI = viper.GetString("database.uri")
	cfg.Logger.Level = viper.GetString("logger.level")
	return cfg
}
