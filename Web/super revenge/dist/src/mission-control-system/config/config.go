package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database     DatabaseConfig
	Server       ServerConfig
	DefaultAdmin DefaultAdminConfig
	JWTSecret    string
	Env          string
	Flag         string
}

type DefaultAdminConfig struct {
	Email    string
	Password string
	Name     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Env      string
}

type ServerConfig struct {
	Port string
}

func (d *DatabaseConfig) GetDatabaseString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func LoadConfig() (*Config, error) {

	godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "rock.roll"),
			Name:     getEnv("DB_NAME", "todo_api"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		DefaultAdmin: DefaultAdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@example.com"),
			Password: getEnv("ADMIN_PASSWORD", "admin"),
			Name:     getEnv("ADMIN_NAME", "Admin"),
		},
		Flag: getEnv("FLAG", "FLAG_NOT_SET"),
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		JWTSecret: getEnv("JWT_SECRET", "your-256-bit-secret"),
		Env:       getEnv("ENV", "development"),
	}

	fmt.Printf("Config %+v", cfg)
	return cfg, nil
}

func getEnv(key string, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Global config instance
var config *Config

// GetConfig returns the application configuration
// It loads the configuration if it hasn't been loaded yet
func GetConfig() *Config {
	if config == nil {
		var err error
		config, err = LoadConfig()
		if err != nil {
			panic("Failed to load configuration: " + err.Error())
		}
	}
	return config
}
