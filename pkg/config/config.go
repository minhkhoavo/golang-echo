package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "root")
	viper.SetDefault("database.name", "golang")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.env", "development")

	// Enable reading from .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Read from .env file (ignore error if file doesn't exist)
	viper.ReadInConfig()

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Bind environment variables with prefix
	viper.SetEnvPrefix("")

	// Map environment variable names
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.env", "SERVER_ENV")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	log.Println("✓ Configuration loaded successfully")
	return &config, nil
}

// GetDSN returns PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// IsDevelopment returns true if server is in development environment
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if server is in production environment
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}