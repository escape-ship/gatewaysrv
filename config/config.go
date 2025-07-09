package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Services Services `mapstructure:"services"`
	Auth     Auth     `mapstructure:"auth"`
	CORS     CORS     `mapstructure:"cors"`
}

type App struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type Services struct {
	Account Service `mapstructure:"account"`
	Product Service `mapstructure:"product"`
	Payment Service `mapstructure:"payment"`
	Order   Service `mapstructure:"order"`
}

type Service struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Auth struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type CORS struct {
	AllowedOrigins     []string `mapstructure:"allowed_origins"`
	AllowedMethods     []string `mapstructure:"allowed_methods"`
	AllowedHeaders     []string `mapstructure:"allowed_headers"`
	AllowedCredentials bool     `mapstructure:"allowed_credentials"`
	MaxAge             int      `mapstructure:"max_age"`
}

func New() (*Config, error) {
	cfg := &Config{}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Set up Viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(wd)
	v.AddConfigPath(filepath.Join(wd, "config"))
	v.AddConfigPath(".")

	// Enable environment variable support
	v.SetEnvPrefix("GATEWAY")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		slog.Warn("Failed to read config file, using defaults and environment variables", "error", err)
	} else {
		slog.Info("Using config file", "path", v.ConfigFileUsed())
	}

	// Unmarshal configuration
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if c.App.Port <= 0 {
		return fmt.Errorf("app.port must be positive")
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}

	services := map[string]Service{
		"account": c.Services.Account,
		"product": c.Services.Product,
		"payment": c.Services.Payment,
		"order":   c.Services.Order,
	}

	for name, service := range services {
		if service.Host == "" {
			return fmt.Errorf("services.%s.host is required", name)
		}
		if service.Port <= 0 {
			return fmt.Errorf("services.%s.port must be positive", name)
		}
	}

	return nil
}

func (c *Config) GetServiceAddress(serviceName string) string {
	var service Service

	switch serviceName {
	case "account":
		service = c.Services.Account
	case "product":
		service = c.Services.Product
	case "payment":
		service = c.Services.Payment
	case "order":
		service = c.Services.Order
	default:
		return ""
	}

	return fmt.Sprintf("%s:%d", service.Host, service.Port)
}
