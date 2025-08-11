package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig              `mapstructure:"server"`
	Database  DatabaseConfig            `mapstructure:"database"`
	Databases map[string]DatabaseConfig `mapstructure:"databases"`
	Sentry    SentryConfig              `mapstructure:"sentry"`
	Metrics   MetricsConfig             `mapstructure:"metrics"`
	Logger    LoggerConfig              `mapstructure:"logger"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type DatabaseConfig struct {
	URI string `mapstructure:"uri"`
}

type SentryConfig struct {
	DSN         string `mapstructure:"dsn"`
	Environment string `mapstructure:"environment"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

func Load() *Config {
	// Read config.yaml if present
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	// Merge .env if present
	viper.SetConfigFile(".env")
	_ = viper.MergeInConfig()

	// Allow ENV variables to override, replace "." in keys with "_"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("sentry.environment", "development")
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("logger.level", "info")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	// Backward compatibility: if Databases is empty but Database.URI provided, seed "main"
	if len(config.Databases) == 0 && config.Database.URI != "" {
		config.Databases = map[string]DatabaseConfig{
			"main": {URI: config.Database.URI},
		}
	}

	return &config
}
