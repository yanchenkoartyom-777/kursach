package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Session   SessionConfig   `yaml:"session"`
	Database  DatabaseConfig  `yaml:"database"`
	Email     EmailConfig     `yaml:"email"`
	TwoFactor TwoFactorConfig `yaml:"two_factor"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	TemplatesPath string `yaml:"templates_path"`
	StaticPath    string `yaml:"static_path"`
}

type SessionConfig struct {
	Secret string `yaml:"secret"`
	MaxAge int    `yaml:"max_age"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type EmailConfig struct {
	SMTPHost    string `yaml:"smtp_host"`
	SMTPPort    int    `yaml:"smtp_port"`
	FromEmail   string `yaml:"from_email"`
	AppPassword string `yaml:"app_password"`
}

type TwoFactorConfig struct {
	Enabled        bool   `yaml:"enabled"`
	CodeLength     int    `yaml:"code_length"`
	CodeExpiration string `yaml:"code_expiration"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if sessionSecret := os.Getenv("SESSION_SECRET"); sessionSecret != "" {
		cfg.Session.Secret = sessionSecret
	}
	if emailPassword := os.Getenv("EMAIL_APP_PASSWORD"); emailPassword != "" {
		cfg.Email.AppPassword = emailPassword
	}
	if emailFrom := os.Getenv("EMAIL_FROM"); emailFrom != "" {
		cfg.Email.FromEmail = emailFrom
	}

	return &cfg, nil
}

// GetCodeExpirationDuration returns the code expiration as time.Duration
func (c *Config) GetCodeExpirationDuration() (time.Duration, error) {
	return time.ParseDuration(c.TwoFactor.CodeExpiration)
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
