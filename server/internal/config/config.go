package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Temporal TemporalConfig `mapstructure:"temporal"`
	// Add other configuration sections as needed (e.g., Cache, Logging)
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Address string `mapstructure:"address"` // e.g., ":8080"
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"` // Data Source Name
}

// TemporalConfig holds Temporal client configuration.
type TemporalConfig struct {
	HostPort  string `mapstructure:"hostPort"` // e.g., "localhost:7233"
	Namespace string `mapstructure:"namespace"`
	TaskQueue string `mapstructure:"taskQueue"`
}

// LoadConfig reads configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	var cfg Config

	// Set default values
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("database.dsn", "root:@tcp(127.0.0.1:4000)/link?charset=utf8mb4&parseTime=True&loc=Local") // Default for local TiDB via docker-compose
	viper.SetDefault("temporal.hostPort", "localhost:7233")
	viper.SetDefault("temporal.namespace", "default")
	viper.SetDefault("temporal.taskQueue", "link-service")

	// Load from config file (e.g., config.yaml) if it exists
	viper.SetConfigName("config")       // Name of config file (without extension)
	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")            // Look for config in the working directory
	viper.AddConfigPath("./server")     // Look for config in the server directory
	viper.AddConfigPath("/etc/link/")   // Path to look for the config file in
	viper.AddConfigPath("$HOME/.link") // Call multiple times to add many search paths

	// Attempt to read the config file but ignore errors if it doesn't exist
	_ = viper.ReadInConfig()

	// Enable environment variable overriding
	// Replaces "." with "_" and converts to uppercase (e.g., server.address becomes SERVER_ADDRESS)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal the config into the Config struct
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
} 